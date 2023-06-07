"""ssl.star 
example struct containing openssl-related data

From https://itnext.io/practical-guide-to-securing-grpc-connections-with-go-and-tls-part-1-f63058e9d6d1
"""

# openssl genrsa -out ca.key 4096
ca_key = """
-----BEGIN RSA PRIVATE KEY-----
MIIJKQIBAAKCAgEAy1oQt++d9oviz1J6G0Xp7HqPHBtUvAUJlQylzbAR/DkbZ3K9
ZRRt+XCPcNhZuTq1KKB6prUqLr0xX6NajsXKn7G523IE17qYKfcL1E3YOUTaouFM
JKpnJAgLGYi92kz6yji4PIww6OkuyKpQP5UkiOdIo8YnttWB3AAzEC7g8HWe8c38
9piBMuc6mayy2Xy9SsW1453be1Iy1dBeudftFFl9L44ZMSfaNObJMg3O09bPhAKm
bMRhwwQ9iWPZUGW66ty8JWed+KuzbfaBQXfB31BaiqaZqw3Wv8nYO6WHINmaNIHR
4BrwZOL6seufrl4LnpTfzMiMDJt38W3bXf0WjB5qZF7QUXvxZ8Pn+sZb+x03YSA3
j+axXeDE5to/fAsGhOCyn+Fao6WCeV/v20a5Cn41TNC0BfjASbETdKz54ksayvnc
+lAM2pqdF/2+ZzqUN7uzRckK5ac99qWP47P/Ye1jhnKgB9D9UfbDQSVZCjO2mJ2l
UqvFTyvLiVgYbY+pwZ98exwocGFbyP2IrM7Ev3kNrCwUqDUbQsZhuNEeMarm1G1d
q+KlBqLt8wEPESls/Pmw1wCbMajOPVjIlBtddpIXs3Zw/+iGvjicjIh47NdZQez/
X6/C7Kw8y2GP9MUOoX+knLhaw8oOdqStanzE3JoJMemQARwfVyrA0f+Q97cCAwEA
AQKCAgBfycosSqQXGee6DzjTlghNy6GT9M/iTWEpI68Kh9DBBcmB3kuWzJvNLxdy
aYdCOIRTYdzEoHwBTj9utI0YydTbiqVo2HmtgQjiY6vf0tdyipuOtB/g+Z/iGiPY
YFBF/5L3JOasJsF3RTgzb/6jJMbz8jaGZvYYKtSj5DgpfFubCVzYvFZXdpkNeFxj
PTV2O0sTaLR4Rsi3e43Up/WnBy53MnxEpWP6grJHzxqhCF4P7ZUMsw7gF1WRvnKa
QD2CoJj+vwGlgPypwX+g4cgbJaVeYwRzYWzrZXZuG09PMXbIo5f++dP6A9aPP1gl
7T7nrQc+KRSO6z0FR0qloEEAMhKnEB5q2YRPk7fLdBUrAYo0goQCIVHnt06Alw/5
SO4q+jM8iOAMtmYKEvRIbHa+YKcdlhZLsJAevd9MxdzFQZwTZP4sBs1hddQA/4cI
tBv9Lh0bYlgWGj75vuyoIJIRADyxDLyR86uNGWawZyV47v87uKXoqLMhuWAKetak
JMG5fuIO1TCtl1+2u3LVq/CoW8cPTpx7JwbGrsw7wj5XBoepy+x8bNywgrZ10mBm
pHZw744AJ68DMs5cOc73vEApKjtkvuDW9WbZucnvUf25o4Y1A0J8gRVl5K46YbIX
jZPlHvULKd7Zy8cAPIZPHMaQEJa/HWqXroeLyY5DEjFleSTUiQKCAQEA/Vw2waHY
cR7QgakiiduO27NvHnHQQ96f7dG9RaUKu2ry0eJw5mLjeDFIyGHHStuHGkIyTx5G
/lpbu0IML8cLG0sIbIiglTdiJsbOrEFOSYHF/opgG9G13WQHupsj3zSzxFRrhnnn
HLWFgXxlERQYqCrEZrqjZPqEQN+ty/28C37aNV0z5WWxiykVrHVF8If4Lg+HJCTV
FpP0j6FemIPuvVocf7vlrJR/DTuhVYPAhdh4w4eGpQ4lV2ko1QVIOPj7GqzBwTUE
DCl+U321tIR36w8ndNB8lXl4OooQnfZV05AQOjYREwkLjU1ARkjS7VxvdDjV9U25
eegYP1nchxH0HQKCAQEAzXh23wGzUF6NLQ+kgnsdftFvH15MZpNrtxbiSj+SyPe8
g9EKg6GlHlAb0BvkAg5FhrJMhzV5mj7VPzxl362trG+ZYTFrhiF3EC9YwgV7o5eY
ED038+q7ECTjsBGiy+E3lvJwuUUT2eOXxLqdcySr0nekbXvNGcgywhzRLaLp1Zq3
DvlqhfPHT3SkWbjizAgg1KBXTf45FpfelWboEQ7Rsj9FcCf0CniDmO3Gz+XstIn6
n7MWcZzlG/j1AmhsR7kd0PKjCugNu/nttTsA+h5c6YHfdoziNNjXWauv5S5ooNKp
VCSJfAg9NFERD0rz5Y+emE/XEI0gPjtai+8oqh/q4wKCAQAj0VKIY0oHC/UsL24L
kTeMBbzyz+JChgmUBG++lcuDnWYAmAOf/mDsEAObGH+lLI3X/32/Q6eDs+B+A6NX
acs/K4dgWJxjG/ZLRxXWslDQAYGtL4DQzf/o8YhKMD6NApVbbxfYZglvPJZILP7Y
wD+QHqOvZjlNQEFMLpMSYKeh9GgC3U9F4e+Mnd1LiTS/AWnrkRRo3rAlRftwBr8p
zpUEveWDhVu93yxrAYAYZ8zi3yyLb/BwCyTqS5qTKvD/5OsS5VNq7gTJd2A9i2sR
vxx45aaNVCAYvZhqpjQdMMMHarlwkU4uo7u3WTF5/jebiNLU2mgdCsTq4A31fs23
ZqldAoIBAQDLwKK4YHpNv4Vl4vYzAh1srgjw5VUD/zq4s/OwxzwrGCgT207+22Pf
HHeINrAzLa3adaMYDXpJ7/cNnzoyxorLzVsfG5/RwgvMu/bbaA6EWobLy7lZozLf
PoWfCs4SOYMjp8UKpCqcTmopBxmtnfbZXhVrEHKCF5nmDieMhto1HRhcvA7bSLQj
4bo80u/sfj74Owx4Zhp8ghuSshp9F+HwTXfxUV1aqMlu9JPLg+jn20/x3+jovzof
NBDa02xU74hWtNXjsdw0xRHpPtqoLUXbtRNA/1IuL73VyUBDF3Nfz6dkrlq76Xuw
DpfJP31+7p3J0pqlah1IORmAXKhJlB8bAoIBAQD3b8qEekhpcU7zSfm1ooEixdHp
w3fUYel4iBp4pMTMaQpGoBMbW5Aa6+XfJRWGVOGgbv+Pzk3dHbDfxFoq7bLGnj41
TRfxVFt69DeAgvn1iiJUjZyRKoRZwktf3aZngiISSPFmlQr7f1a+iuMOkYShjm+R
+WZpNsBUJNdWOW3ky7wPUevfkbRY4NhVZjqjyy6JxMSphS4sWyuLfHWunDyjbCbh
MY/F9J98AsXeDjJWg/Vsqd5BWIP27qpidc2XAC4IiqpmG08hify6ofNvg5PLZdFp
LpLUOIIDFOpzPjv485h1qDDke/TqsgZBoNuEpkgHdCcSVc+4yNaqR8MBSG2Y
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

# certificate.conf
# @see https://gist.github.com/yidas/af42d2952d85c0951c1722fcd68716c6
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

# openssl req -new -key service.key -out service.csr -config certificate.conf
service_csr = """
-----BEGIN CERTIFICATE REQUEST-----
MIIFFjCCAv4CAQAwgYwxCzAJBgNVBAYTAlRXMQ8wDQYDVQQIDAZUYWl3YW4xDzAN
BgNVBAcMBlRhaXBlaTEOMAwGA1UECgwFWUlEQVMxEDAOBgNVBAsMB1NlcnZpY2Ux
IDAeBgkqhkiG9w0BCQEWEXlvdXJtYWlsQG1haWwuY29tMRcwFQYDVQQDDA55b3Vy
ZG9tYWluLmNvbTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAKsFdAWs
ndW7Rm64IUaXVdqmNgkzOhR9Wkm9KEX7LRpNpSaJ8D4oJDpxOLwrCHql/pGkaNse
jt0+NjdYdDpvHRabQ2hSC3klzSK+DkclBl3wt1PRe3r82wxuH5rN7fYxKPNIyc3/
AzIy4B+zDWke1qvCosDX9Jx1w2YhH+xLYP+ebOWvd1QqCp8/UUDQ+ajT1JZ6g9iy
5geVhXsGntWZMEuTSlXorK/o0DlGFUi1CwhFe3SUgJtDc3GWVOADOXhnZl+Eh1W2
p0PITzVfBAH+L29/A4axFULSADD8Yl7BnrcnEXSLaPqv2+ypg1NOoFN0qqiwR1sW
2lt3fPonFwEngL5YPeL7lLv4NabzRqOHg0drcbNI8AUlB1IGPuaLyXWo+cGQdmVD
FqdEQFZdulS0SwU66QM6PUGJcCW2IaAhfl8h/Zu8HIQZIRFuWyX/XY23XeEGgVFz
ZAE3COijlglTzuVGsxujYxrLhVmUb9Lj+LysSgGxjcwSw/f6HZJlCHSMH7Y2PGxI
lbrqf6CbPT4DonLN15yRe68h7KYFhMAmV+6P4QRjDmYALrQbKdk1ObNaZhFbrHdM
z8ZRJe2dWTx4cTSUlclC/51DeBsAoFZaEYgYdn03w9PWyc3IUgaffL7eXSX/nEs7
sj3QM2/MkpCX+vQwLV/kQuL9+Hwvk4RZhAcnAgMBAAGgRDBCBgkqhkiG9w0BCQ4x
NTAzMDEGA1UdEQQqMCiCECoueW91cmRvbWFpbi5jb22CFCouZGV2LnlvdXJkb21h
aW4uY29tMA0GCSqGSIb3DQEBCwUAA4ICAQA4hb/ze700Fn0k6kSpHKO4exeh7rci
s2e5mYsool9XskLyZMvZR1ZQpb7N3NtMTAxrH73nCDXLQZXQSKC8xPvV3lE2034j
tiCD5eShpN3PzLJXQigKcm/TTnfq0iz99tWRZ+O6UbrBgw+Mz9l1z2cv5nL/rAVv
p+nbi5UIZl9KmwBoAHJ+02YMRg6aV4kRkmYwFCTsZ/EjxYLF2CNUqLaEi+bTyoEl
glvY9NWxcYjSWe7xrIe0zG9DYASBtI8Rkh9Vq9/DWyghyMbg8rSb793Giq9zKsqu
CyxV1XQtIiPbdHPkiH0IC2NpfBq5dx8TFv+poqSF43RlDxY7FF+OJ3yY8e5pH9GC
7pcmGoJMrjJAme/RrC2akvJi4h8v3bTTpa3DwU7mdPCx4bjve/XsTa2C2lNmSJLO
7Y8ek3w/WS38p7EkNtrTnWXogLA6fcuMEx7NRzA2+Nk/5tI0MJSiNsBufmg7lP9t
8dhfxknPNGkpTIGYViFLDd5IYDvYq+kD+cWcdsYJMvw/wkM33Wl/VvVcoqAy5VTd
LOMstpDnKluX/zl3nEzmPr2SF0KVDssz2jkr+WtmcUYI4To5hJLjOzOy1akiRjbk
0flhuHRlo67fsHZBHRdQLrhP19poFsnXHdQEQ685MXHOJz8BgrFau7L6suFZLQu2
i+i8f3NQsg7XiA==
-----END CERTIFICATE REQUEST-----
"""

# openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext
service_pem = """
-----BEGIN CERTIFICATE-----
MIIFcjCCA1qgAwIBAgIJAJJQ3VGH9o16MA0GCSqGSIb3DQEBCwUAMC0xCzAJBgNV
BAYTAlVTMQswCQYDVQQIDAJOSjERMA8GA1UECgwIQ0EsIEluYy4wHhcNMjMwNjA2
MDMyNjUzWhcNMjQwNjA1MDMyNjUzWjCBjDELMAkGA1UEBhMCVFcxDzANBgNVBAgM
BlRhaXdhbjEPMA0GA1UEBwwGVGFpcGVpMQ4wDAYDVQQKDAVZSURBUzEQMA4GA1UE
CwwHU2VydmljZTEgMB4GCSqGSIb3DQEJARYReW91cm1haWxAbWFpbC5jb20xFzAV
BgNVBAMMDnlvdXJkb21haW4uY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIIC
CgKCAgEAqwV0Bayd1btGbrghRpdV2qY2CTM6FH1aSb0oRfstGk2lJonwPigkOnE4
vCsIeqX+kaRo2x6O3T42N1h0Om8dFptDaFILeSXNIr4ORyUGXfC3U9F7evzbDG4f
ms3t9jEo80jJzf8DMjLgH7MNaR7Wq8KiwNf0nHXDZiEf7Etg/55s5a93VCoKnz9R
QND5qNPUlnqD2LLmB5WFewae1ZkwS5NKVeisr+jQOUYVSLULCEV7dJSAm0NzcZZU
4AM5eGdmX4SHVbanQ8hPNV8EAf4vb38DhrEVQtIAMPxiXsGetycRdIto+q/b7KmD
U06gU3SqqLBHWxbaW3d8+icXASeAvlg94vuUu/g1pvNGo4eDR2txs0jwBSUHUgY+
5ovJdaj5wZB2ZUMWp0RAVl26VLRLBTrpAzo9QYlwJbYhoCF+XyH9m7wchBkhEW5b
Jf9djbdd4QaBUXNkATcI6KOWCVPO5UazG6NjGsuFWZRv0uP4vKxKAbGNzBLD9/od
kmUIdIwftjY8bEiVuup/oJs9PgOics3XnJF7ryHspgWEwCZX7o/hBGMOZgAutBsp
2TU5s1pmEVusd0zPxlEl7Z1ZPHhxNJSVyUL/nUN4GwCgVloRiBh2fTfD09bJzchS
Bp98vt5dJf+cSzuyPdAzb8ySkJf69DAtX+RC4v34fC+ThFmEBycCAwEAAaM1MDMw
MQYDVR0RBCowKIIQKi55b3VyZG9tYWluLmNvbYIUKi5kZXYueW91cmRvbWFpbi5j
b20wDQYJKoZIhvcNAQELBQADggIBALYgLFAxhMBEkxIdsDyYsOrXBvP3I37HuyIo
rVQ1+7chI0ymi0Rsrf7VP76AlMb7xVm37elPwKOh0rBpda2hP2/JV7pCOXXv2Nwf
e930NeJXqL6wDSPQZye8gnvECVHqovdFIawe5s1MnD3srqi/jElhMV1iIO16EatM
Y4ayASM1sFIFjzAWm+eoSYuusGEyMnLoVlrzgappEMy8IneUosODgfpIeUdoaofK
BpVobegJJBUE5XnuwnWlkkmJHHcH70bTG1znG7+5g1I/fBHOV3dmvz3JFxhVBz1f
gebWM6NyIrTJG6pQl313s/3Qmoe2T2T2CjvI5Kc653QeI8idT3XP+w/vGS8hKTMk
7qSn8obwlFmjQ0Po68rRULwAco8RDNDlIL/5J/2uyJnR6aH290cI8QQP1jVT7DZ5
bZ13Z2behKIogHcDv70bfsUuniu2D5Nr5uw7s/rY0L/awAEtyMuMuhHGLe94Q+q0
DTj1r8NuMLsCk1rnT6wTuBM3Bb5mFfquxQIJDSWeICMWQLAz+YDtdpoOQETaruUD
uiJvqE5b6HnZFIjCZsZ7NxePg5Kk7cmeve3rZ4ipKzy00OXMF/t03kgL2uJYUrrJ
UsdNljxvOrAk96iU4KZb1YiyLKJiH68z1btjH5UWU1bCxHr3y7JttgFEx/BY/Hsp
YYoLVe9y
-----END CERTIFICATE-----
"""

ca_bundle = struct(
    ca_cert = ca_cert,
    ca_key = ca_key,
    service_csr = service_csr,
    service_key = service_key,
    service_pem = service_pem,
)
