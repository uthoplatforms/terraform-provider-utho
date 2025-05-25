resource "utho_domain" "example" {
  domain = "example.com"
}

resource "utho_dns_record" "example" {
  domain   = utho_domain.example.domain
  type     = "A"
  hostname = "subdomain.${utho_domain.example.domain}"
  value    = "1.1.1.1"
  ttl      = "65444"
  porttype = "TCP"
  port     = "5060"
  priority = "10"
  weight   = "100"
}
