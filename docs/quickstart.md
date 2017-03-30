Quickstart
==========

Summary
-------

- Prepare your AWS account:
	- Create [S3 bucket](https://aws.amazon.com/s3/)
	- Create IAM credentials for S3
- Set up Writer ("the server")
- Run an example application feeding off of Pyramid

Pyramid runs just fine whether you have the Writer and application on the same server or not.


Create S3 bucket
----------------

I created an S3 bucket named `eventhorizon1.fn61.net` in region `EU (Frankfurt)`.

So our details are:

- Bucket name = `eventhorizon1.fn61.net` (it doesn't have to be a DNS name)
- Region code = `eu-central-1` ([S3 region codes](http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region))

You'll need these later.


Create IAM credentials for S3
-----------------------------

[Create AWS IAM credentials](configuring/create-aws-iam-credentials.md)


Assembling your ENV variable
----------------------------

Pyramid is configured via one ENV variable: `STORE`. It contains the S3 details
and all the rest is handled automatically.

The variable looks like this:

```
STORE=s3://AKIAIZQ7QCQOTMYODIAA:secretkey@eu-central-1/eventhorizon1.fn61.net
```

Where:

```
STORE=s3://APIKEY_ID:APIKEY_SECRET@S3_REGION/S3_BUCKET
```

NOTE: temporarily you have to replace `/` chars in secret key with `_`.


Install & configure Writer on the server
----------------------------------------

Find out the **public IP address** of your Writer server:

```
$ ifconfig eth0
eth0      Link encap:Ethernet  HWaddr 06:2B:12:10:B3:0B
          inet addr:1.2.3.4
```

[Enter Pyramid CLI](enter-pyramid-cli.md) and bootstrap the Writer cluster:

```
$ pyramid writer-bootstrap
2017/03/25 16:35:30 writer-bootstrap: Usage: <WriterIp>
$ pyramid writer-bootstrap 1.2.3.4
2017/03/25 16:35:51 bootstrap: generating certificate authority
2017/03/25 16:35:51 bootstrap: generating auth token
2017/03/25 16:35:51 bootstrap: generating discovery file
2017/03/25 16:35:51 bootstrap: uploading discovery file to scalablestore
2017/03/25 16:35:51 bootstrap: bootstrapped Writer cluster for ip 1.2.3.4
```

The above command essentially uploads a JSON file to your S3 bucket to let clients
know how to connect to Writer servers.

Now start the Writer on the server:

```
$ docker run --name pyramid -d --net=host -e "STORE=$STORE" fn61/pyramid

then check the logs:

$ docker logs pyramid
2017/03/22 14:30:58        .
2017/03/22 14:30:58       /=\\       PyramidDB
2017/03/22 14:30:58      /===\ \     function61.com
2017/03/22 14:30:58     /=====\  \
2017/03/22 14:30:58    /=======\  /
2017/03/22 14:30:58   /=========\/
2017/03/22 14:30:58 configfactory: downloading discovery file
2017/03/22 14:30:58 PubSubServer: binding to 0.0.0.0:9091
2017/03/22 14:30:58 CompressedEncryptedStore: mkdir /pyramid-data/store-compressed_and_encrypted
...
```

Everything seems ok. Now enter Pyramid CLI again to poke with the system.

We'll now create the minimum two streams required to run the system. Create root stream:

```
$ pyramid stream-create /
```

Create `/_sub` stream (subscriptions will reside as sub-streams under this path):

```
$ pyramid stream-create /_sub
```

Ok now the Writer has been properly set up!


Now run your example application
--------------------------------

Go run [the example app](https://github.com/function61/pyramid-exampleapp-go),
either on the same server or a different server.
