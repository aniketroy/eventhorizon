package config

const (
	WALMANAGER_DATADIR = "/eventhorizon/store-live"

	LONGTERMSHIPPER_PATH = "/eventhorizon/store-longterm"

	SEEKABLE_STORE_PATH = "/eventhorizon/store-seekable"

	COMPRESSED_ENCRYPTED_STORE_PATH = "/eventhorizon/store-compressed_and_encrypted"

	BOLTDB_DIR = "/eventhorizon"

	S3_BUCKET = "eventhorizon.fn61.net"

	S3_BUCKET_REGION = "us-east-1"

	PUBSUB_PORT = 9091

	WRITER_HTTP_PORT = 9092

	WAL_SIZE_THRESHOLD = uint64(4 * 1024 * 1024)

	CHUNK_ROTATE_THRESHOLD = 8 * 1024 * 1024
)
