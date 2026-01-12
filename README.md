# Dağıtık ve Hata-Toleranslı Disk Kayıt Sistemi (HaToKuSe)

![Go](https://img.shields.io/badge/Language-Go-blue)
![gRPC](https://img.shields.io/badge/Protocol-gRPC-green)
![Architecture](https://img.shields.io/badge/Architecture-Leader--Member-orange)

Bu proje, **Sistem Programlama** dersi ödevi kapsamında geliştirilmiş; dağıtık, hata tolere edebilen (fault-tolerant) ve dinamik ölçeklenebilen bir disk kayıt sistemidir. Proje, "HaToKuSe" (Hata-Tolere Kuyruk Servisi) protokolü temel alınarak **Go** programlama dili ile gerçekleştirilmiştir.

---

## 📋 İçindekiler

- [Proje Hakkında](#-proje-hakkında)
- [Sistem Mimarisi](#-sistem-mimarisi)
- [Özellikler ve Gereksinimler](#-özellikler-ve-gereksinimler)
- [Kurulum ve Çalıştırma](#-kurulum-ve-çalıştırma)
- [Kullanım (İstemci Bağlantısı)](#-kullanım-istemci-bağlantısı)
- [Yapılandırma (Tolerance)](#-yapılandırma-tolerance)
- [Test Senaryoları](#-test-senaryoları)
- [Proje Yapısı](#-proje-yapısı)

---

## 📖 Proje Hakkında

Bu projenin temel amacı, merkezi bir Lider ve dinamik Aile Üyelerinden (Members) oluşan dağıtık bir yapı kurmaktır.

Sistem, soket üzerinde HTTP veya SMTP gibi hazır protokoller yerine, ödev kapsamında belirtilen **ilkel metin tabanlı protokol (HaToKuSe)** ve **gRPC** ile haberleşmektedir. İstemciler (Clients) sisteme dışarıdan bağlanarak sadece Lider ile muhatap olurken, veriler arka planda belirlenen hata tolerans seviyesine göre üyelere yedekli bir şekilde dağıtılır.

---

## 🏗 Sistem Mimarisi

Sistem **Lider-Üye (Leader-Member)** mimarisine dayanmaktadır.

1.  **Lider (Leader):**
    - Dış dünyadan (İstemciden) gelen `SET` ve `GET` isteklerini karşılar.
    - `tolerance.conf` dosyasındaki değere göre (N), veriyi N adet üyeye **gRPC** üzerinden kopyalar.
    - Hangi verinin hangi üyede olduğunu bellekte tutar ve yük dengelemesini (Load Balancing) sağlar.
2.  **Aile Üyesi (Member):**
    - Liderden gelen veriyi diske `.msg` formatında yazar (Persistent Storage).
    - İstendiğinde diski okuyup veriyi Lidere döner.
    - Dinamik olarak sisteme katılıp (Join), sistemden ayrılabilir (Crash/Leave).

_(Not: Proje kapsamında özel bir İstemci (Client) yazılımı geliştirilmemiş olup, sistem `netcat`, `telnet` veya ödev kapsamında sağlanan test araçları ile haberleşecek şekilde tasarlanmıştır.)_

---

## ✨ Özellikler

- **İletişim Protokolleri:**
  - Client ↔ Leader: **TCP Soket** (Metin tabanlı: `SET`, `GET`)
  - Leader ↔ Member: **gRPC** (Protobuf nesneleri ile)
- **Veri Saklama:** Mesajlar disk üzerinde `messages/` klasöründe kalıcı olarak saklanır.
- **Hata Toleransı (Fault Tolerance):** Yapılandırılabilir (`tolerance.conf`) yedekleme sayısı (1-7 arası).
- **Dinamik Üyelik:** Yeni üyeler sisteme çalışma zamanında (runtime) katılabilir.
- **Crash Recovery:** Veriyi tutan bir üye çökse bile, Lider diğer kopyalardan veriyi getirir.

---

## 🚀 Kurulum ve Çalıştırma

Proje tek bir `main` dosyası üzerinden çalışır. Uygulama başlatıldığında port durumuna göre **Lider** mi yoksa **Üye** mi olacağına kendisi karar verir.

### 1. Derleme (Build)

Projeyi derleyerek çalıştırılabilir bir `main` dosyası oluşturun:

```bash
go build -o main cmd/main.go
```

### 2. Çalıştırma (Run)

Sistemi ayağa kaldırmak için terminalden `main` dosyasını çalıştırmanız yeterlidir veya derlemeden doğrudan run ile çalıştırılabilir.

```bash
go run cmd/main.go
```

**Terminal 1 - Lideri Başlat:**
İlk çalışan süreç 5555 portunu alır ve LİDER olur.

```bash
./main
# Çıktı: [LIDER] 5555 portunda dinliyor...
```

**Terminal 2, 3... - Üyeleri Başlat:**
Lider ayaktayken çalıştırılan diğer süreçler, sıradaki boş portları alarak ÜYE (Follower) olur ve kümeye katılır.

```bash
./main
# Çıktı: [UYE] 5556 portunda çalışıyor, Lidere bağlandı.
```

---

## 💻 Kullanım (İstemci Bağlantısı)

Sistemi test etmek için herhangi bir TCP istemcisi (`nc`, `telnet` vb.) kullanabilirsiniz. Lider varsayılan olarak **TCP 6666** portundan komut bekler.

**Bağlantı:**

```bash
telnet 127.0.0.1 6666
```

### Komut Seti

| Komut              | Açıklama                      | Örnek                   |
| ------------------ | ----------------------------- | ----------------------- |
| `SET <id> <mesaj>` | Mesajı sisteme kaydeder.      | `SET 100 Merhaba Dunya` |
| `GET <id>`         | ID'si verilen mesajı getirir. | `GET 100`               |

**Örnek Senaryo:**

```text
> SET 42 SistemProgramlama
< OK

> GET 42
< SistemProgramlama

> GET 999
< NOT_FOUND
```

---

## ⚙️ Yapılandırma

Kök dizindeki `tolerance.conf` dosyası ile hata tolerans seviyesini belirleyebilirsiniz:

```ini
tolerance=2
```

Bu ayar, bir verinin başarılı sayılması için en az kaç farklı üyeye (node) yazılması gerektiğini belirtir.

---

## 🧪 Test Senaryoları

Proje geliştirilirken aşağıdaki senaryolar başarıyla test edilmiştir:

1.  **Yük Dağılımı (Load Balancing):**
    - `tolerance=2` ve 4 üye ile yapılan testte, 1000 adet mesajın üyelere yaklaşık eşit (~500'er) dağıtıldığı gözlemlenmiştir.
2.  **Crash & Recovery (Hata Kurtarma):**
    - Bir veriyi tutan üyelerden biri (process kill ile) kapatıldığında, Liderin bu durumu fark edip veriyi hayatta kalan diğer üyeden çekebildiği doğrulanmıştır.

---

## 🛠 Proje Geliştirme Adımları

Bu bölüm, projenin başlangıçtan üretim seviyesine kadar izlenen geliştirme adımlarını özetler. Her adım, kapsam ve teknik uygulama detaylarıyla birlikte verilmiştir (kişi isimleri intentionally belirtilmemiştir).

1. **Proje İskeletinin Kurulumu:**

   - Go modül başlatma (`go mod init`), temel klasör yapısının oluşturulması (`cmd/`, `internal/`, `proto/`, `messages/`).
   - Bağımlılık yönetimi ve minimal çalışma örneği ile ilk derleme/doğrulama.

2. **Üyelik (Family) gRPC Protokolünün Tasarımı:**

   - `proto/family.proto` içinde `NodeInfo`, `FamilyView`, `StoredMessage`, `MessageId` mesajları ve `FamilyService` servisinin tanımlanması.
   - RPC uçları: `Join`, `GetFamily`, `Store`, `Retrieve` fonksiyonlarının sözleşmelerinin belirlenmesi.

3. **Protobuf’tan Go Kodlarının Üretilmesi:**

   - `protoc` ve Go eklentileri ile (`protoc-gen-go`, `protoc-gen-go-grpc`) kod üretimi.
   - `option go_package` ile doğru paket adlandırması ve import uyumluluğu.

4. **Lider için Dinamik gRPC Üyelik Mantığı:**

   - Lider tarafında `Registry` yapısının kullanımı; üyelerin `Join` ile kümeye katılması ve `GetFamily` ile görünümün yayılması.
   - Üyelerin çalışma zamanında eklenmesi/çıkması senaryolarının ele alınması.

5. **TCP Komut Ayrıştırıcı (6666 Portu):**

   - Basit metin tabanlı protokol (`SET <id> <text>`, `GET <id>`) için ayrıştırıcı.
   - Hata yönetimi ve geçersiz komutların geri bildirimleri.

6. **Disk Kalıcılık (Single Node Storage):**

   - `internal/storage/disk.go` ile `messages/` altında `.msg` dosyalarına yazma/okuma.
   - `bufio` kullanarak güvenli yazma/okuma ve dosya senkronizasyonu (`Sync`).

7. **gRPC Tabanlı StorageService’ın Uygulanması:**

   - Üye düğümünde `Store`/`Retrieve` RPC’lerinin implementasyonu.
   - Liderin üyelere istemci olarak bağlanıp (client stub) replikasyon için çağrıları gerçekleştirmesi.

8. **Tolerance Değerinin Yüklenmesi:**

   - `tolerance.conf` üzerinden yapılandırma okuma (`internal/config/config.go`).
   - Format/doğrulama ve hatalı/eksik durumlar için anlamlı hata mesajları.

9. **LeaderCoordinator ile SET/GET ve Replikasyon Akışı:**

   - Liderin gelen `SET` isteğini yerelde saklayıp `tolerance` kadar üyeye gRPC ile kopyalaması.
   - `GET` isteğinde yerel okuma başarısızsa üyelere sorgu göndererek veriyi getirmesi.

10. **Round-Robin Üye Seçimi:**

    - `nextIndex` ve dairesel gezinme ile eşit dağılımlı replikasyon.
    - Kendi düğümünün ve yinelenen seçimlerin dışlanması; `tolerance` kadar seçimin garanti edilmesi.

11. **Crash Yönetimi ve Recovery Mekanizmaları:**

    - Üyeye bağlantı hatasında üyenin “ulaşılamaz” durumunun algılanması ve bir sonraki üyeden deneme.
    - Ayrıntılı günlükler (attempt/skip/fail) ile fault-tolerance davranışının izlenmesi.

12. **Loglama ve Gözlemleme:**
    - Periyodik aile görünümü ve mesaj replikasyon durumunun yazdırılması (ör. `StartFamilyPrinter`, `StartMessagePrinter`).
    - Operasyonel izlenebilirlik ve hata ayıklama için yalın ama bilgilendirici log formatları.

---

## 📂 Proje Yapısı

```
distributed-disk-register-with-grpc/
├── cmd/
│   └── main.go              # Ana uygulama dosyası (Role Detection)
├── internal/
│   ├── common/              # Ortak yapılar
│   ├── config/              # Konfigürasyon okuma
│   ├── leader/              # Lider mantığı (TCP server, Load Balancer)
│   ├── node/                # Üye mantığı (Disk storage, gRPC server)
│   ├── storage/             # Disk I/O işlemleri
│   └── discovery/           # Port keşfetme servisi
├── proto/
│   ├── family.proto         # gRPC protokol tanımları
│   └── family/              # Generate edilmiş Go kodları
├── messages/                # Verilerin kaydedildiği klasör
├── tolerance.conf           # Ayar dosyası
└── README.md
```

---

**Takım Üyeleri:**

- Ali Ellikci
- Kadir Kopuz
- Burak Aktürk
- Mehmet Akyürek
