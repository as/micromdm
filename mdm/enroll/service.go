package enroll

import (
	"bytes"
	"crypto/x509"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/groob/plist"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/as/micromdm/platform/config"
	"github.com/as/micromdm/platform/profile"
	"github.com/as/micromdm/platform/pubsub"
)

const (
	EnrollmentProfileId string = "com.github.micromdm.micromdm.enroll"
	OTAProfileId        string = "com.github.micromdm.micromdm.ota"
)

type Service interface {
	Enroll(ctx context.Context) (profile.Mobileconfig, error)
	OTAEnroll(ctx context.Context) (profile.Mobileconfig, error)
	OTAPhase2(ctx context.Context) (profile.Mobileconfig, error)
	OTAPhase3(ctx context.Context) (profile.Mobileconfig, error)
}

func NewService(topic TopicProvider, sub pubsub.Subscriber, caCertPath, scepURL, scepChallenge, url, tlsCertPath, scepSubject string, profileDB profile.Store) (Service, error) {
	var caCert, tlsCert []byte
	var err error

	if caCertPath != "" {
		caCert, err = ioutil.ReadFile(caCertPath)

		if err != nil {
			return nil, err
		}
	}

	if tlsCertPath != "" {
		tlsCert, err = ioutil.ReadFile(tlsCertPath)

		if err != nil {
			return nil, err
		}
	}

	if scepSubject == "" {
		scepSubject = "/O=MicroMDM/CN=MicroMDM Identity (%ComputerName%)"
	}

	subjectElements := strings.Split(scepSubject, "/")
	var subject [][][]string

	for _, element := range subjectElements {
		if element == "" {
			continue
		}
		subjectKeyValue := strings.Split(element, "=")
		subject = append(subject, [][]string{[]string{subjectKeyValue[0], subjectKeyValue[1]}})
	}

	// fetch the push topic from the db.
	// will be "" if the push certificate hasn't been uploaded yet
	pushTopic, _ := topic.PushTopic()
	svc := &service{
		URL:           url,
		SCEPURL:       scepURL,
		SCEPSubject:   subject,
		SCEPChallenge: scepChallenge,
		CACert:        caCert,
		TLSCert:       tlsCert,
		ProfileDB:     profileDB,
		Topic:         pushTopic,
		topicProvier:  topic,
	}

	if err := updateTopic(svc, sub); err != nil {
		return nil, errors.Wrap(err, "enroll: start topic update goroutine")
	}

	return svc, nil
}

func updateTopic(svc *service, sub pubsub.Subscriber) error {
	configEvents, err := sub.Subscribe(context.TODO(), "enroll-server-configs", config.ConfigTopic)
	if err != nil {
		return errors.Wrap(err, "update enrollment service")
	}
	go func() {
		for {
			select {
			case <-configEvents:
				topic, err := svc.topicProvier.PushTopic()
				if err != nil {
					log.Println("enroll: get push topic %s", topic)
				}
				svc.mu.Lock()	// TODO(as): fix
				svc.Topic = topic
				svc.mu.Unlock()

				// terminate the loop here because the topic should never change
				goto exit
			}
		}
	exit:
		return
	}()
	return nil
}

type service struct {
	URL           string
	SCEPURL       string
	SCEPChallenge string
	SCEPSubject   [][][]string
	CACert        []byte
	TLSCert       []byte
	ProfileDB     profile.Store

	topicProvier TopicProvider

	mu    sync.RWMutex
	Topic string // APNS Topic for MDM notifications
}

type TopicProvider interface {
	PushTopic() (string, error)
}

func profileOrPayloadFromFunc(f interface{}) (interface{}, error) {
	fPayload, ok := f.(func() (Payload, error))
	if !ok {
		fProfile := f.(func() (Profile, error))
		return fProfile()
	}
	return fPayload()
}

func profileOrPayloadToMobileconfig(in interface{}) (profile.Mobileconfig, error) {
	if _, ok := in.(Payload); !ok {
		_ = in.(Profile)
	}
	buf := new(bytes.Buffer)
	enc := plist.NewEncoder(buf)
	enc.Indent("  ")
	err := enc.Encode(in)
	return buf.Bytes(), err
}

func (svc *service) findOrMakeMobileconfig(id string, f interface{}) (profile.Mobileconfig, error) {
	p, err := svc.ProfileDB.ProfileById(id)
	if err != nil {
		if profile.IsNotFound(err) {
			profile, err := profileOrPayloadFromFunc(f)
			if err != nil {
				return nil, err
			}
			return profileOrPayloadToMobileconfig(profile)
		}
		return nil, err
	}
	return p.Mobileconfig, nil
}

func (svc *service) Enroll(ctx context.Context) (profile.Mobileconfig, error) {
	return svc.findOrMakeMobileconfig(EnrollmentProfileId, svc.MakeEnrollmentProfile)
}

const perUserConnections = "com.apple.mdm.per-user-connections"

func (svc *service) MakeEnrollmentProfile() (Profile, error) {
	profile := NewProfile()
	profile.PayloadIdentifier = EnrollmentProfileId
	profile.PayloadOrganization = "MicroMDM"
	profile.PayloadDisplayName = "Enrollment Profile"
	profile.PayloadDescription = "The server may alter your settings"
	profile.PayloadScope = "System"

	mdmPayload := NewPayload("com.apple.mdm")
	mdmPayload.PayloadDescription = "Enrolls with the MDM server"
	mdmPayload.PayloadOrganization = "MicroMDM"
	mdmPayload.PayloadIdentifier = EnrollmentProfileId + ".mdm"
	mdmPayload.PayloadScope = "System"

	svc.mu.Lock()
	topic := svc.Topic
	svc.mu.Unlock()

	mdmPayloadContent := MDMPayloadContent{
		Payload:             *mdmPayload,
		AccessRights:        allRights(),
		CheckInURL:          svc.URL + "/mdm/checkin",
		CheckOutWhenRemoved: true,
		ServerURL:           svc.URL + "/mdm/connect",
		Topic:               topic,
		SignMessage:         true,
		ServerCapabilities:  []string{perUserConnections},
	}

	payloadContent := []interface{}{}

	if svc.SCEPURL != "" {
		scepContent := SCEPPayloadContent{
			URL:      svc.SCEPURL,
			Keysize:  2048,
			KeyType:  "RSA",
			KeyUsage: int(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment),
			Name:     "Device Management Identity Certificate",
			Subject:  svc.SCEPSubject,
		}

		if svc.SCEPChallenge != "" {
			scepContent.Challenge = svc.SCEPChallenge
		}

		scepPayload := NewPayload("com.apple.security.scep")
		scepPayload.PayloadDescription = "Configures SCEP"
		scepPayload.PayloadDisplayName = "SCEP"
		scepPayload.PayloadIdentifier = EnrollmentProfileId + ".scep"
		scepPayload.PayloadOrganization = "MicroMDM"
		scepPayload.PayloadContent = scepContent
		scepPayload.PayloadScope = "System"

		payloadContent = append(payloadContent, *scepPayload)
		mdmPayloadContent.IdentityCertificateUUID = scepPayload.PayloadUUID
	}

	payloadContent = append(payloadContent, mdmPayloadContent)

	if len(svc.CACert) > 0 {
		caPayload := NewPayload("com.apple.security.root")
		caPayload.PayloadDisplayName = "Root certificate for MicroMDM"
		caPayload.PayloadDescription = "Installs the root CA certificate for MicroMDM"
		caPayload.PayloadIdentifier = EnrollmentProfileId + ".cert.ca"
		caPayload.PayloadContent = svc.CACert

		payloadContent = append(payloadContent, *caPayload)
	}

	// Client needs to trust us at this point if we are using a self signed certificate.
	if len(svc.TLSCert) > 0 {
		tlsPayload := NewPayload("com.apple.security.pem")
		tlsPayload.PayloadDisplayName = "Self-signed TLS certificate for MicroMDM"
		tlsPayload.PayloadDescription = "Installs the TLS certificate for MicroMDM"
		tlsPayload.PayloadIdentifier = EnrollmentProfileId + ".cert.selfsigned"
		tlsPayload.PayloadContent = svc.TLSCert

		payloadContent = append(payloadContent, *tlsPayload)
	}

	profile.PayloadContent = payloadContent

	return *profile, nil
}

// OTAEnroll returns an Over-the-Air "Profile Service" Payload for enrollment.
func (svc *service) OTAEnroll(ctx context.Context) (profile.Mobileconfig, error) {
	return svc.findOrMakeMobileconfig(OTAProfileId, svc.MakeOTAEnrollPayload)
}

func (svc *service) MakeOTAEnrollPayload() (Payload, error) {
	payload := NewPayload("Profile Service")
	payload.PayloadIdentifier = OTAProfileId
	payload.PayloadDisplayName = "MicroMDM Profile Service"
	payload.PayloadDescription = "Profile Service enrollment"
	payload.PayloadOrganization = "MicroMDM"
	payload.PayloadContent = ProfileServicePayload{
		URL:              svc.URL + "/ota/phase23",
		Challenge:        "",
		DeviceAttributes: []string{"UDID", "VERSION", "PRODUCT", "SERIAL", "MEID", "IMEI"},
	}

	// yes, this is a bare Payload, not a Profile
	return *payload, nil
}

// OTAPhase2 returns a SCEP Profile for use in phase 2 of Over-the-Air enrollment.
func (svc *service) OTAPhase2(ctx context.Context) (profile.Mobileconfig, error) {
	return svc.findOrMakeMobileconfig(OTAProfileId+".phase2", svc.MakeOTAPhase2Profile)
}

func (svc *service) MakeOTAPhase2Profile() (Profile, error) {
	profile := NewProfile()
	profile.PayloadIdentifier = OTAProfileId + ".phase2"
	profile.PayloadOrganization = "MicroMDM"
	profile.PayloadDisplayName = "OTA Phase 2"
	profile.PayloadDescription = "The server may alter your settings"
	profile.PayloadScope = "System"

	scepContent := SCEPPayloadContent{
		URL:      svc.SCEPURL,
		Keysize:  2048, // NOTE: OTA docs recommend 1024
		KeyType:  "RSA",
		KeyUsage: int(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment),
		Name:     "OTA Phase 2 Certificate",
		Subject:  svc.SCEPSubject,
	}

	if svc.SCEPChallenge != "" {
		scepContent.Challenge = svc.SCEPChallenge
	}

	scepPayload := NewPayload("com.apple.security.scep")
	scepPayload.PayloadDescription = "Configures SCEP"
	scepPayload.PayloadDisplayName = "SCEP"
	scepPayload.PayloadIdentifier = OTAProfileId + ".phase2.scep"
	scepPayload.PayloadOrganization = "MicroMDM"
	scepPayload.PayloadContent = scepContent
	scepPayload.PayloadScope = "System"

	profile.PayloadContent = append(profile.PayloadContent, *scepPayload)

	return *profile, nil
}

// OTAPhase3 returns a Profile for use in phase 3 of Over-the-Air profile enrollment.
// This would typically be the final or end profile of the Over-the-Air
// enrollment process. In our case this would probably be a device-specifc
// MDM enrollment payload.
// TODO: Not implemented.
func (svc *service) OTAPhase3(ctx context.Context) (profile.Mobileconfig, error) {
	return profile.Mobileconfig{}, nil
}
