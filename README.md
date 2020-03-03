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

## References

* https://github.com/smancke/mailck
* https://github.com/go-graphite/carbonapi/tree/master/zipper