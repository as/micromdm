package builtin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/as/micromdm/pkg/crypto"
	"github.com/as/micromdm/platform/config"
	"github.com/as/micromdm/platform/pubsub"
)

const (
	ConfigBucket = "mdm.ServerConfig"
)

// DB stores server configuration in BoltDB
type DB struct {
	*bolt.DB
	Publisher pubsub.Publisher
}

func NewDB(db *bolt.DB, pub pubsub.Publisher) (*DB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ConfigBucket))
		return err
	})
	store := &DB{DB: db, Publisher: pub}
	return store, err
}

func (db *DB) SavePushCertificate(cert, key []byte) error {
	tx, err := db.DB.Begin(true)
	if err != nil {
		return errors.Wrap(err, "begin transaction to store push certificate in bolt")
	}
	bkt := tx.Bucket([]byte(ConfigBucket))
	if bkt == nil {
		return fmt.Errorf("config: bucket %q not found", ConfigBucket)
	}
	pb, err := config.MarshalServerConfig(&config.ServerConfig{
		PushCertificate: cert,
		PrivateKey:      key,
	})
	if err != nil {
		return errors.Wrap(err, "save push cert in bolt bucket")
	}

	if err := bkt.Put([]byte("config"), pb); err != nil {
		return errors.Wrap(err, "save ServerConfig in bucket")
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	if err := db.Publisher.Publish(context.TODO(), config.ConfigTopic, []byte("updated")); err != nil {
		return err
	}
	return err
}

func (db *DB) serverConfig() (*config.ServerConfig, error) {
	var conf config.ServerConfig
	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(ConfigBucket))
		data := bkt.Get([]byte("config"))
		if data == nil {
			return &notFound{"ServerConfig", "no config found in boltdb"}
		}
		return config.UnmarshalServerConfig(data, &conf)
	})
	return &conf, errors.Wrap(err, "get server config from bolt")
}

func (db *DB) PushCertificate() (*tls.Certificate, error) {
	conf, err := db.serverConfig()
	if err != nil {
		return nil, errors.Wrap(err, "get server config for push cert")
	}

	// load private key
	pkeyBlock, _ := pem.Decode(conf.PrivateKey)
	if pkeyBlock == nil {
		return nil, errors.New("decode private key for push cert")
	}

	priv, err := x509.ParsePKCS1PrivateKey(pkeyBlock.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse push certificate key from server config")
	}

	// load certificate
	certBlock, _ := pem.Decode(conf.PushCertificate)
	if certBlock == nil {
		return nil, errors.New("decode push certificate PEM")
	}

	pushCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse push certificate from server config")
	}

	cert := tls.Certificate{
		Certificate: [][]byte{pushCert.Raw},
		PrivateKey:  priv,
		Leaf:        pushCert,
	}
	return &cert, nil
}

func (db *DB) PushTopic() (string, error) {
	cert, err := db.PushCertificate()
	if err != nil {
		return "", errors.Wrap(err, "get push certificate for topic")
	}
	topic, err := crypto.TopicFromCert(cert.Leaf)
	return topic, errors.Wrap(err, "get topic from push certificate")
}

type notFound struct {
	ResourceType string
	Message      string
}

func (e *notFound) Error() string {
	return fmt.Sprintf("not found: %s %s", e.ResourceType, e.Message)
}
