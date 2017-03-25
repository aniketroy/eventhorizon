
Summary
-------

Enables log-based architectures with realtime streaming, unlimited storage, compression and encryption.

| Key facts               |                                                                                          |
|-------------------------|------------------------------------------------------------------------------------------|
| Technologies            | Golang, Linux, Docker (not required for use), AWS S3.                                    |
| Production ready?       | No, not yet. But not too far away.                                                       |
| Msg delivery semantics  | Exactly once, exact in-order delivery within a stream.                                   |
| Use cases               | Eventsourcing, high availability software/datacenter failover, microservices, big data.  |
| Storage capacity        | Practically unlimited. You have to pay your bills though. :)                             |
| Data durability         | All writes transactionally backed by Write-Ahead-Log just like in databases.             |
| High availability       | Planned, going to use Hashicorp's Raft implementation.                                   |
| Data stored at          | AWS S3. Google Storage support planned.                                                  |
| Encryption at transport | TLS (CA & server certs automatically managed)                                            |
| Encryption at rest      | AES-CBC. Encryption keys are not trusted to AWS (=> not using S3 server-side encryption) |
| Security                | [Our security policy & information](https://function61.com/security/)                    |


Docs
----

- [Quickstart](docs/quickstart.md) (for users)
- [Example app using Pyramid](https://github.com/function61/pyramid-exampleapp-go)
  (best way to grasp what Pyramid is about, but you should also read quickstart)
- [Building & contributing](docs/building-and-contributing.md) (for developers)
- [Operating](docs/operating.md)
- [Security](docs/architecture/security.md)
- [Network](docs/architecture/network.md)
- [TODO & roadmap](docs/todo-roadmap.md)
- [Alternatives](docs/alternatives.md)


Architecture summary
--------------------

Writer manages writes to streams. Streams are divided into ~8 MB chunks which
are stored compressed & encrypted in AWS S3, except the last, "live", chunk that
we're writing into.

![](docs/architecture/diagram.png)

Data flow at a glance:

- Somebody contacts Writer to append a line to the stream
- Writer notifies pub/sub that a stream has new lines
- All interested applications that want to read from any streams have a Pusher,
  which subscribes to pub/sub notifications for those streams.
- Once the Pusher learns of new events, it issues a read to either:
	- a Writer if it's live data OR
	- S3 if it's older data.
- Pusher then pushes this read to your application. Your application transactionally
  verifies that the read offset is at the offset of last stored write in database.
  If not, error is returned along with the correct offset, and Pusher continues
  to read from the correct offset.

Programming language support:

- You only have to implement the tiny [pushlib](pusher/pushlib) for your
  programming language (look at the architecture diagram), which essentially is
  just one HTTP path with incoming JSON data.
- pushlib is already implemented for these programming languages:
	- Go ([example app](https://github.com/function61/pyramid-exampleapp-go))


License & fair play
-------------------

Legal license: Apache 2.0 - free software.

Moral license: if you benefit commercially from the use of this project, any help
would be appreciated (though not legally required):

- Order [support contract/consulting from us](https://function61.com/consulting/)
- Contribute with new features, bug fixes and/or help with issues
- Become a sponsor
- Donate money
