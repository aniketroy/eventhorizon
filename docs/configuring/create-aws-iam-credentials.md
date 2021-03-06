Create AWS IAM credentials
==========================

Make IAM credentials for S3
---------------------------

We need to set up two users in IAM, both with their own auth keys:

- eventhorizon-writer
	- has read-write privileges
- eventhorizon-pusher
	- has read-only privileges

Create them in IAM with "programmatic access" checkbox, so we get API keys.

Do not attach any policies to them, as we'll set up a super restrictive custom policy.

We just got two sets of access keys generated:

- eventhorizon-writer = AKIAIZQ7QCQOTMYODIAA secret = ...
- eventhorizon-pusher = AKIAJPNFVKAFJQHRBY3A secret = ...

(of course API key IDs will be different for you)


Set up custom policies for both accounts
----------------------------------------

Now, go to `IAM > eventhorizon-writer > add inline policy > custom policy`:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": ["s3:ListBucket"],
            "Resource": ["arn:aws:s3:::eventhorizon1.fn61.net"]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObject"
            ],
            "Resource": ["arn:aws:s3:::eventhorizon1.fn61.net/*"]
        }
    ]
}
```

Policy name: s3-eventhorizon-readwrite

And now go to `IAM > eventhorizon-pusher > add inline policy > custom policy`:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": ["s3:ListBucket"],
            "Resource": ["arn:aws:s3:::eventhorizon1.fn61.net"]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject"
            ],
            "Resource": [
                "arn:aws:s3:::eventhorizon1.fn61.net/*"
            ]
        }
    ]
}
```

Policy name: s3-eventhorizon-read.
