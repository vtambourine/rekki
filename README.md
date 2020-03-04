# rekki

## Usage

Pull and run image:

```sh
docker run -t -p 8080:8080 vtambourine/rekki
```

Once the image is up, send *POST* requests to the `/email/validate` endpoint:

```sh
$ curl -XPOST -d '{"email":"xxx@yyy.zzz"}' http://localhost:8080/email/validate
{
  "valid": false,
  "validators": {
    "blacklist": {
      "valid": true
    },
    "hostname": {
      "valid": false,
      "reason": "INVALID_TLD"
    },
    "regexp": {
      "valid": true
    },
    "smtp": {
      "valid": false,
      "reason": "INVALID_HOSTNAME"
    }
  }
}

$ curl -XPOST -d '{"email":"support@rekki.com"}' http://localhost:8080/email/validate
{
  "valid": true,
  "validators": {
    "blacklist": {
      "valid": true
    },
    "hostname": {
      "valid": true
    },
    "regexp": {
      "valid": true
    },
    "smtp": {
      "valid": true
    }
  }
}

$ curl -XPOST -d '{"email":"mailboxthatdoesntexist@gmail.com"}' http://localhost:8080/email/validate
{
  "valid": false,
  "validators": {
    "blacklist": {
      "valid": true
    },
    "hostname": {
      "valid": true
    },
    "regexp": {
      "valid": true
    },
    "smtp": {
      "valid": false,
      "reason": "UNAVAILABLE_MAILBOX"
    }
  }
}
```

## Limitations and Improvements

As any educational or technical assignment project, this project contains many shortcuts and improvements points. Here are some of the most obvious improvements and missing features:

* Unit tests.
* Request logging.
* Confugaration of default parameters, such as application port, logging levels, SMPT port and from email address.
* Blacklist and top-level domains lists are served as static files loaded into memory. Better solution would be to use some memory cache, managable without applicaion restart.
* Scripts for updating domains lists from respected sources.
* Better regexp validation with support of internationalized domains and UTF characters.
* More validators, such as domain reputation and TXT records.
* Current implementation runs all available validators independently. Quick improvemnet would be to introduce validators order and relation. For instance, once blacklist validation fails, do not run SMPT validation checks.

## References

* [smancke/mailck](https://github.com/smancke/mailck)
* [go-graphite/carbonapi/zipper](https://github.com/go-graphite/carbonapi/tree/master/zipper)
