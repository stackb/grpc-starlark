"""tls.grpc.star 

Demonstrates the sending and revcing of a message with mutual TLS.

"""
pb = proto.package("example.routeguide")
service_name = "example.routeguide.RouteGuide"

# certificate.conf
certificate_conf = """
[req]
default_bits = 2048
distinguished_name = dn
prompt             = no
req_extensions = req_ext

[dn]
C="TW"
ST="Taiwan"
L="Taipei"
O="YIDAS"
OU="Service"
emailAddress="yourmail@mail.com"
CN="yourdomain.com"

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.0 = *.yourdomain.com
DNS.1 = *.dev.yourdomain.com
"""

# openssl genrsa -out ca.key 4096
ca_key = """
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
"""

# openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
ca_cert = """
-----BEGIN CERTIFICATE-----
MIIE1jCCAr4CCQDCCUzEk7A4HTANBgkqhkiG9w0BAQsFADAtMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCTkoxETAPBgNVBAoMCENBLCBJbmMuMB4XDTIzMDYwNjAzMTIx
OVoXDTI0MDYwNTAzMTIxOVowLTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAk5KMREw
DwYDVQQKDAhDQSwgSW5jLjCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIB
AMtaELfvnfaL4s9SehtF6ex6jxwbVLwFCZUMpc2wEfw5G2dyvWUUbflwj3DYWbk6
tSigeqa1Ki69MV+jWo7Fyp+xudtyBNe6mCn3C9RN2DlE2qLhTCSqZyQICxmIvdpM
+so4uDyMMOjpLsiqUD+VJIjnSKPGJ7bVgdwAMxAu4PB1nvHN/PaYgTLnOpmsstl8
vUrFteOd23tSMtXQXrnX7RRZfS+OGTEn2jTmyTINztPWz4QCpmzEYcMEPYlj2VBl
uurcvCVnnfirs232gUF3wd9QWoqmmasN1r/J2DulhyDZmjSB0eAa8GTi+rHrn65e
C56U38zIjAybd/Ft2139FoweamRe0FF78WfD5/rGW/sdN2EgN4/msV3gxObaP3wL
BoTgsp/hWqOlgnlf79tGuQp+NUzQtAX4wEmxE3Ss+eJLGsr53PpQDNqanRf9vmc6
lDe7s0XJCuWnPfalj+Oz/2HtY4ZyoAfQ/VH2w0ElWQoztpidpVKrxU8ry4lYGG2P
qcGffHscKHBhW8j9iKzOxL95DawsFKg1G0LGYbjRHjGq5tRtXavipQai7fMBDxEp
bPz5sNcAmzGozj1YyJQbXXaSF7N2cP/ohr44nIyIeOzXWUHs/1+vwuysPMthj/TF
DqF/pJy4WsPKDnakrWp8xNyaCTHpkAEcH1cqwNH/kPe3AgMBAAEwDQYJKoZIhvcN
AQELBQADggIBAMP24mIJKl+mJnS9LEYr6Wni0atHw3vVreBfrTxz+UmKyNNrJd2v
ULNDJZfLh3jAsPwcZMzgdYaaYr1bG/1Vp7+sdrX8wseWlx/lHITqyxbP2DctvDMh
yjEhcfuXNirlyU3LMJZS8vfX9DwMoqgYptyeDXvkiL49+fy80Fua95zFmxVCumFD
/fq43uKs1tirw1zBlOTQtmhlu/BAPZV/iujesVfcgBE8vd8xIWNQB37ghAypIH41
gKnYKTKFkVRRCjRZzDB39ApdjxMwTjER6H0oZYrwShiWDZHm8nCdBPCk5wC46WVF
mNtly8XCD4jr1X11o6jtO5lc21G5AzOjYbr5VS9gwtAKXoZVJkivLae6O9QsTeKW
M9OrPLpfBwiVin2O2NqNq391vyGe8rggtkJnkC+8l2lOtZ1+htXVyX99LV+0JakA
kq3sUbT3YSXQrec4kj6QAcu5TGg/erFktripXPv25xtc9DhR0GaZhW+04UnmRkz/
cS2LxtKsMMic+m9sEgdTOMZOJ233SajBBndATR+pFMYfTYvxskD88PtJsXAl7oYj
zFY+ETMIpJ4HtXBWDcOqQ6oHTqn10QFvWRIz5mWaMxqd+h1d+XVVlE0GxJXJYPew
dR0HRbQR0RLv+toj0XPP8GlnvDbhWX40UDDzWqEZjJKs3ZfIM4tOyamG
-----END CERTIFICATE-----
"""

# openssl genrsa -out service.key 4096
service_key = """
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
"""

# openssl req -new -key service.key -out service.csr -config certificate.conf
service_csr = """
"""

root_ca_certificate = """
"""

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

root_cas = crypto.x509.CertPool()
root_cas.append(root_ca_certificate)
print("root_cas:", root_cas)

config = crypto.tls.Config(
    certificates = [certificate],
    client_auth = crypto.tls.ClientAuthType.NONE,
    insecure_skip_verify = True,
    root_certificate_authorities = root_cas,
)
print("config:", config)

print("insecure transport credentials:", grpc.credentials.Insecure())
creds = grpc.credentials.Tls(config)
print("secure transport credentials:", creds)

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
