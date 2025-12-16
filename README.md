# Distributed Disk Register with gRPC

Bu proje, Sistem Programlama dersi kapsamında istenen
**dağıtık, hata-tolere mesaj kayıt sisteminin**
Go (Golang) ve gRPC kullanılarak geliştirilmiş halidir.

## 📌 Proje Amacı
- Lider–üye mimarisi ile çalışan
- Hata toleransı destekleyen
- Mesajları disk üzerinde saklayan
- Dinamik üye katılımına izin veren
bir dağıtık kayıt sistemi geliştirmek.

## 🧱 Sistem Mimarisi
- **Client**
  - Lider sunucuya text tabanlı `SET` ve `GET` istekleri gönderir.
- **Leader**
  - Client’tan gelen mesajları alır
  - `tolerance.conf` dosyasına göre mesajları üyelere dağıtır
  - Hata toleransı sağlandıktan sonra client’a `OK / ERROR` döner
- **Member**
  - Liderden gRPC üzerinden gelen mesajları alır
  - Mesajları disk üzerinde saklar
  - Periyodik olarak tuttuğu mesaj sayısını raporlar

## 🔌 İletişim
- Client ↔ Leader : Text tabanlı protokol
- Leader ↔ Member : gRPC (.protobuf)

## 📁 Proje Dizini
client/ → İstemci uygulaması
leader/ → Lider sunucu
member/ → Aile üyesi sunucular
proto/ → gRPC protobuf tanımları
config/ → tolerance.conf vb. ayarlar
internal/ → Ortak yardımcı kodlar

markdown
Kodu kopyala

## 👥 Takım Çalışması
- Geliştirme süreci GitHub Projects üzerinden yürütülmektedir.
- Her özellik ayrı bir task ve feature branch olarak geliştirilmektedir.
- Tamamlanan işler Pull Request ile `main` branch’e merge edilmektedir.

## ⚙️ Kullanılan Teknolojiler
- Go (Golang)
- gRPC
- Protocol Buffers
- Git & GitHub Projects