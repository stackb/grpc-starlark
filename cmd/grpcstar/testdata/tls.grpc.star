"""headers.grpc.star 

Demonstrates the sending and revcing of a message with mutual TLS.

"""
pb = proto.package("example.routeguide")
service_name = "example.routeguide.RouteGuide"

# openssl genrsa -out server.key 2048
private_key = """
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA/IWmKuQm/1py64Ieq1YlqbR3XNLLNAfnU6b4dHwyLLEYiCmk
7BKYp6NtLrpUTFqnFMoWLBsWASkLOr0bsJBd5azV4eWis/B4I4MzAJ5NnVJkOVTu
NHFTEsib2pD7SouMtLVVe2/29AbGj6dYWzoqLoytklOvI6q633Hc/dT4p2QVn01N
S4UCoAyGtSNQvZXIXO3pKf5BG51DqYOzBE9BnRoUGN3WOGG6y0wFw8X80nCI2ncv
pDPBfn+QtT3YAFENBBtGUZyqOUn2tZgJS0fCfT4/ykHO6zQFvheMsRsKlXznuUtL
bqtpMlL0Fdtjrtv6ZZlOwu69mchWqkmFioZo3QIDAQABAoIBADc0JHJl9ByIsmzH
wlqkd5FU8W8qad/TBoAkFVapu/JHONyzdelh21tyf7DibQFQJAyIbTZxKWtRhLHv
m3kK5mwKT6uVnu8FV84zpVeyQ7drxps99OEkEQwfLOsoHLdcMINkzO4yOON6A7ht
1gQDgCsy99LwVm5OqZGle7FF+KHm8SRkukz1v8S3MTjP3S5BzSBneXbNzZeg5Ied
68P4ghzJ0MtHprgDx4r1bnueGVYpIWQgehwjr0qYbVZoQ3d6XHEO7GOxx9ba8lWa
d7/BWynIR3ABR3MSB4tpeyODWa3mp817QWyMGNudDz9Pl+wbodQf5Xz0mRcjvKiC
q0DTK4UCgYEA/sZWYaJAs6AQbTUigIeH2uPprwxP34k2dW+JN3L4rFTcSxc/aJzw
4eqVi7EcKAutdu/CPUMzKbeqIA+H8pI8ju/5AuE3x07efcp1Em4jOEblg4X+EAgZ
CCHVbjcRCYLVUDlIc9kvf0BuJIukV6UDRnjcm5Gg1mw5yMfEF32aJ6sCgYEA/byJ
zc8KJuUhGylkyXRqUWy2W+r73n/vVbskCaQYWJ0inWfjr5ncB437lXJEQBbSX5Fw
WI6bJYDm6p6kTFIVWgNHB00MJH/EdBD8PL418fQZyABH1FKquwr4XdH15hri0BYx
kc++UsQQG7TDLzVvdZYdv/0uRnTtGW6svpecCZcCgYBHDT8n6V0L+za5jhj6KVH8
/JS+KbvYxmZ2p81ntludi+kH1Art/N680nQ0SgdlL6SHx+OuvB/3oW4DlPE/+AKF
hm02nWK15cvs3tp5clfGKRd275ZkGC4K84yXOSo6Mc+VmPQYwtgZL/nHnV4Ox0k7
jRdRF3L4eaQ/115bgr7MEwKBgQCYhxG/qknL/8ja7xMrFtQihltI/gTSR82zl3+e
XApWmn8IaD8yfCcMU4l82Oe2LwHfeSoz0eXpsYceWqchSeaT6Yx1ExfNiRCrRNqc
GSuMetRUqfaD5/3B2mJa47AR1u+pbu31XRBn6HxWa185rcGGyeqwUp3StM8ijqlB
GRovmQKBgEQpI0xH+uS4HeRzfgNjVKGg9Pus/XoJ2HObpyoYp9Mucn2GPuYWGXGv
f44OozrJXMj0lbM1ijL0mZojrZmGErTu9nbrnAycyBZ44c/2VoXupH+JgnjoHfeM
Ei1Wpk8lmqQM3BZS1Mv4sxKA9itJr7FVE/zjr4VS1FkdyVIesXLh
-----END RSA PRIVATE KEY-----
"""

# openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
public_key = """
-----BEGIN CERTIFICATE-----
MIIDSjCCAjICCQCwp/mPa5xmijANBgkqhkiG9w0BAQsFADBnMQswCQYDVQQGEwJ1
czELMAkGA1UECAwCY28xEDAOBgNVBAcMB2JvdWxkZXIxFDASBgNVBAoMC3N0YWNr
LmJ1aWxkMQ0wCwYDVQQLDARncnBjMRQwEgYDVQQDDAtzdGFjay5idWlsZDAeFw0y
MzA2MDYwMDEwMjdaFw0zMzA2MDMwMDEwMjdaMGcxCzAJBgNVBAYTAnVzMQswCQYD
VQQIDAJjbzEQMA4GA1UEBwwHYm91bGRlcjEUMBIGA1UECgwLc3RhY2suYnVpbGQx
DTALBgNVBAsMBGdycGMxFDASBgNVBAMMC3N0YWNrLmJ1aWxkMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEA/IWmKuQm/1py64Ieq1YlqbR3XNLLNAfnU6b4
dHwyLLEYiCmk7BKYp6NtLrpUTFqnFMoWLBsWASkLOr0bsJBd5azV4eWis/B4I4Mz
AJ5NnVJkOVTuNHFTEsib2pD7SouMtLVVe2/29AbGj6dYWzoqLoytklOvI6q633Hc
/dT4p2QVn01NS4UCoAyGtSNQvZXIXO3pKf5BG51DqYOzBE9BnRoUGN3WOGG6y0wF
w8X80nCI2ncvpDPBfn+QtT3YAFENBBtGUZyqOUn2tZgJS0fCfT4/ykHO6zQFvheM
sRsKlXznuUtLbqtpMlL0Fdtjrtv6ZZlOwu69mchWqkmFioZo3QIDAQABMA0GCSqG
SIb3DQEBCwUAA4IBAQAdYUOvbbClQzOTghSQSr9Y3F0EjsQ8Wyji+IjVugUikguv
EVw4qOaUmkPrkwwIF1/PVPdIG1dt7VC40FeQwyvd2kunAXREzMkY8mMXwrNWT3Ls
cc0SaIAchF7U34U0yFPDC5JWEoKPXHZTrsY3+sxhXMEIdhu7Ls1JIRAPa4mEKQiM
qwj27L1hZ/po3HBtMqU0RgM1RMSZRdyRVftTjKZPkNjIV9NGS+7OWp1dBXHDGgN+
xi9VgFHxRifxjQrbhZ6GrQ5eix+86OXw6Hm0yF6Md3vrxPRt3gws36RJYaFHcJim
/q4MADkw+J2qHriBnzA3PttJY6Rx+PCGlKFXQ3ci
-----END CERTIFICATE-----
"""

certificate = crypto.tls.Certificate(public_key, private_key)
print("certificate:", certificate)

config = crypto.tls.Config(
    certificates = [certificate],
    client_auth = crypto.tls.ClientAuthType.NONE,
)
print("config:", config)

creds = grpc.credentials.Tls(config)
print("transport credentials:", creds)

# === [Server Handler Functions] ================================================

def get_feature(stream, point):
    """get_feature implements a unary method handler

    Args:
        stream: the stream object
        point: the requested Point
    Returns:
        a Feature, ideally nearest to the given point.

    """
    return pb.Feature(name = "point (%d,%d)" % (point.longitude, point.latitude))

listener = net.Listener()
server = grpc.Server(
    credentials = creds,
)
server.register(service_name, {
    "GetFeature": get_feature,
})
thread.defer(lambda: server.start(listener))

channel = grpc.Channel(listener.address, credentials = creds)
client = grpc.Client(service_name, channel)

def call_get_feature():
    feature = client.GetFeature(
        request = pb.Point(longitude = 1, latitude = 2),
    )
    print("client: GetFeature response message:", feature)

call_get_feature()

server.stop()
