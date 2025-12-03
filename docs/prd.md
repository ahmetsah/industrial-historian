---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11]
inputDocuments: ['docs/analysis/research/technical-historian-architecture-research-2025-12-02.md']
workflowType: 'prd'
lastStep: 11
project_name: 'historian'
user_name: 'Ahmet'
date: '2025-12-02T12:13:01+03:00'
---

# Product Requirements Document - historian

**Author:** Ahmet
**Date:** 2025-12-02T12:13:01+03:00

## Executive Summary

**Historian**, endüstriyel veri yönetimini yeniden tanımlayan, modern, yüksek performanslı ve ölçeklenebilir bir **Kurumsal Endüstriyel Historian Platformu**dur. AVEVA PI ve GE Proficy gibi geleneksel devlere meydan okuyan Historian, Rust ve NATS JetStream teknolojileri üzerine inşa edilmiş dağıtık mikroservis mimarisiyle, saniyede 50.000'den fazla olayı işleyebilir ve ağ kesintilerinde bile sıfır veri kaybı garantisi sunar.

Sadece geçmiş verileri depolamakla kalmayan Historian, entegre **Dijital İkiz (Digital Twin)** motoruyla fiziksel süreçleri modeller ve geleceğe yönelik kestirimci analizler sunar. ISA 18.2 Alarm Yönetimi ve FDA 21 CFR Part 11 Denetim İzi uyumluluğu ile kritik endüstriyel standartları karşılarken, kademeli depolama mimarisiyle (RAM -> SSD -> S3) maliyetleri optimize eder.

### What Makes This Special

1.  **Geleceği Simüle Eden Güç (Dijital İkiz):** Rakiplerin aksine, Python (GEKKO) tabanlı fiziksel modelleme motoru ile sadece geçmişi değil, geleceği de simüle ederek proaktif karar destek ve kestirimci bakım yetenekleri sunar.
2.  **Ödünsüz Performans ve Güvenilirlik:** Rust ve NATS JetStream tabanlı mimari, >50k/s yazma hızı ve <100ms sorgu süresi sunarken, "Store and Forward" teknolojisi ile ağ kesintilerinde bile veri bütünlüğünü korur.
3.  **Maliyet Etkin Kademeli Depolama:** Gorilla sıkıştırması ve verinin yaşam döngüsüne göre (Sıcak, Ilık, Soğuk) otomatik taşınması, depolama maliyetlerini rakiplerine göre önemli ölçüde düşürür.
4.  **Yerleşik Endüstriyel Uyumluluk:** Ek modüllere gerek kalmadan, kutudan çıktığı haliyle ISA 18.2 ve FDA 21 CFR Part 11 standartlarını karşılayarak regüle endüstriler için hazır çözüm sunar.

## Project Classification

**Technical Type:** Industrial IoT Platform / SaaS / Embedded
**Domain:** Industrial Automation (OT)
**Complexity:** High

## Success Criteria

### User Success

*   **Güven ve Huzur:** Operatörler, ağ kesintilerinde bile "Store and Forward" mekanizmasının çalıştığını bilerek veri kaybı endişesi yaşamazlar.
*   **Anlık İçgörü:** Mühendisler, 1 yıllık trend verisini 100ms'nin altında sorgulayarak (uPlot & Rust Engine) sorunları anında teşhis edebilirler.
*   **Proaktif Karar Alma:** Yöneticiler, Dijital İkiz (Lite) simülasyonları sayesinde olası darboğazları veya arızaları gerçekleşmeden önce fark ederler.

### Business Success

*   **Pazar Penetrasyonu:** İlk 12 ay içinde 5 pilot endüstriyel tesiste başarılı kurulum ve doğrulama.
*   **Maliyet Avantajı:** Rakiplere (AVEVA, GE) kıyasla %50 daha düşük Toplam Sahip Olma Maliyeti (TCO) sunmak (Donanım verimliliği ve lisanslama modeli ile).
*   **Operasyonel Süreklilik:** %99.99 sistem çalışma süresi (Uptime) sağlamak.

### Technical Success

*   **Yüksek Performans:** Tek bir Ingestion düğümünde >50.000 olay/saniye kararlı yazma hızı.
*   **Depolama Verimliliği:** Gorilla sıkıştırması ile veri noktası başına ortalama <2 bayt depolama alanı kullanımı.
*   **Regülasyon Uyumu:** ISA 18.2 (Alarm) ve FDA 21 CFR Part 11 (Audit Trail) standartlarının teknik gereksinimlerini eksiksiz karşılamak.

### Measurable Outcomes

*   **Sorgu Latency:** P99 < 100ms (1 yıllık veri aralığı için).
*   **Veri Kaybı:** %0 (Ağ kesintisi simülasyonlarında).
*   **Simülasyon Doğruluğu:** Dijital İkiz tahminlerinin gerçekleşen verilerle %90+ korelasyon göstermesi.

## Product Scope

### MVP - Minimum Viable Product

*   **Core Historian:** Rust Ingestor, NATS JetStream, Rust Engine (LSM-Tree), Modbus/OPC UA Adaptörleri.
*   **Digital Twin (Lite):** Python (GEKKO) tabanlı temel fiziksel modelleme ve basit anomali tespiti (Tek değişkenli).
*   **Visualization:** React/uPlot tabanlı gerçek zamanlı trend izleme paneli.
*   **Compliance:** Temel ISA 18.2 Alarm Durum Makinesi ve FDA uyumlu Audit Loglama.

### Growth Features (Post-MVP)

*   **Advanced Digital Twin:** Çok değişkenli karmaşık simülasyonlar ve senaryo analizleri.
*   **Tiered Storage Automation:** Verilerin S3/MinIO'ya (Soğuk Katman) otomatik arşivlenmesi ve buradan sorgulanması.
*   **Advanced Reporting:** Özelleştirilebilir raporlar ve dışa aktarma seçenekleri.

### Vision (Future)

*   **Autonomous Operations:** Yapay zeka destekli, insan müdahalesi olmadan süreç parametrelerini optimize eden otonom karar mekanizması.
*   **Federated Historian:** Bulut ve Uç (Edge) noktaları arasında hibrit çalışan, küresel ölçekli federasyon mimarisi.

## User Journeys

### Journey 1: Operatör Mehmet - The Predictive Hero (Kestirimci Kahraman)

**Persona:** Mehmet (45), 20 yıllık saha operatörü. Teknolojiden çok tecrübesine güvenir, karmaşık ekranlardan nefret eder.
**Hedef:** Vardiyasını kazasız ve duruşsuz tamamlamak.

**Senaryo:** Gece vardiyası, saat 03:00. Kritik Besleme Pompası (P-101) mekanik yorgunluk belirtileri gösteriyor.

**Mevcut Durum (Kaos):** Mehmet, SCADA ekranındaki yüzlerce veri arasında kaybolmuştur. Pompa titreşimi artar ancak alarm eşiğini geçmediği için sistem sessizdir. Saat 04:30'da pompa aniden kilitlenir, üretim durur. Panik, stres ve 2 saatlik üretim kaybı yaşanır.

**Historian ile (Kontrol):**
1.  **Sessiz Tespit (03:05):** Arka planda çalışan Dijital İkiz (Sim), P-101'in akım/basınç ilişkisindeki "Model Sapması"nı tespit eder. Henüz fiziksel bir alarm yoktur.
2.  **Akıllı Uyarı (03:06):** Mehmet'in Viz ekranına kırmızı bir alarm değil, sarı bir "Kestirimci İçgörü" düşer: *"P-101 Performans Sapması. Tahmini Arıza: < 4 Saat."*
3.  **Görsel Teyit:** Mehmet bildirime tıklar. uPlot grafiğinde Gerçekleşen Basınç (Mavi) ile İdeal Basınç (Gri Kesikli - Digital Twin) arasındaki makasın açıldığını net bir şekilde görür.
4.  **Proaktif Müdahale (03:10):** Panik yapmadan, hattı durdurmadan Yedek Pompa P-102'yi devreye alır.
5.  **Sonuç:** Sıfır duruş, sıfır stres. Mehmet, bir kriz çıkmadan onu önlemiştir.

### Journey 2: Kalite Uzmanı Zeynep - The Audit Guardian (Denetim Muhafızı)

**Persona:** Zeynep (32), Titiz ve kuralcı bir Kalite Güvence Uzmanı.
**Hedef:** FDA denetimlerinden "Sıfır Bulgusuz" geçmek.

**Senaryo:** Bir ilaç partisinde (Batch #992) kalite sorunu şüphesi vardır. FDA denetçisi, 6 ay önceki o partinin üretim parametrelerini ve kimin müdahale ettiğini sorar.

**Historian ile:**
1.  **Anında Erişim:** Zeynep, Historian Viz arayüzünde "Batch #992"yi aratır. Saniyeler içinde o partiye ait tüm sıcaklık ve basınç trendleri (Engine) ekrana gelir.
2.  **Denetim İzi (Audit Trail):** Zeynep, "Olay Günlüğü"nü açar. O gün saat 14:15'te Operatör Mehmet'in "Karıştırıcı Hızı" set değerini manuel olarak değiştirdiğini görür.
3.  **Veri Bütünlüğü:** Denetçi, "Bu logun değiştirilmediğini kanıtla" der. Zeynep, tek tıkla logun "Zincirleme Hash" (Blockchain benzeri) imzasını gösterir.
4.  **Sonuç:** Denetçi ikna olur. Zeynep, günler sürecek bir kağıt kürek işinden kurtulur ve denetimi başarıyla tamamlar.

### Journey 3: Bakım Teknisyeni Burak - The Planned Maintainer (Planlı Bakımcı)

**Persona:** Burak (28), Genç ve teknolojiye meraklı bakım teknisyeni.
**Hedef:** Arızalara koşmak yerine, bakımları planlı yapmak.

**Senaryo:** Operatör Mehmet'in uyarısı üzerine P-101 pompasına bakmaya gider.

**Historian ile:**
1.  **Hazırlıklı Gitmek:** Sahaya gitmeden önce tabletinden P-101'in son 1 aylık titreşim trendini (Viz) açar. Trendin son 3 gündür yavaşça arttığını görür.
2.  **Doğru Teşhis:** Dijital İkiz'in "Rulman Aşınması Olasılığı %85" notunu görür. Yanına sadece anahtar takımını değil, yedek rulmanı da alır.
3.  **Hızlı Çözüm:** Pompayı açtığında gerçekten rulmanın dağılmak üzere olduğunu görür. Parçayı değiştirir ve pompayı tekrar devreye alır.
4.  **Sonuç:** Deneme-yanılma yok. Tek seferde doğru müdahale.

### Journey 4: Süreç Mühendisi Ayşe - The Efficiency Architect (Verimlilik Mimarı)

**Persona:** Ayşe (35), Veri odaklı Süreç Mühendisi.
**Hedef:** Üretim hattındaki darboğazları bulup verimliliği (OEE) artırmak.

**Senaryo:** Hattın enerji tüketimini analiz etmek ve optimize etmek istiyor.

**Historian ile:**
1.  **Büyük Veri Sorgusu:** Ayşe, son 1 yılın enerji tüketim verisini (milyarlarca satır) sorgular. Rust Engine ve Downsampling sayesinde grafik 1 saniyenin altında yüklenir.
2.  **Korelasyon Analizi:** Enerji tüketimi ile Üretim Hızı grafiklerini üst üste bindirir (Overlay). Belirli bir ürün tipinde enerji tüketiminin %20 arttığını fark eder.
3.  **İyileştirme:** Bu ürün için fırın sıcaklık set değerini optimize eder.
4.  **Sonuç:** Yıllık %5 enerji tasarrufu sağlar. Veriyle konuşarak süreci iyileştirir.

### Journey 5: Fabrika Müdürü Ali Bey - The Insightful Leader (İçgörülü Lider)

**Persona:** Ali Bey (50), Fabrika Müdürü. Detaylarda boğulmak istemez, büyük resmi görmek ister.
**Hedef:** Karlılığı artırmak ve operasyonel riskleri yönetmek.

**Senaryo:** Sabah toplantısı öncesi fabrikanın durumunu görmek istiyor.

**Historian ile:**
1.  **KPI Dashboard:** Tabletinden Historian Dashboard'unu açar. Anlık OEE, Enerji Tüketimi ve Aktif Alarm Sayısı'nı tek ekranda görür.
2.  **Maliyet Analizi:** "Kademeli Depolama Tasarrufu" widget'ına bakar. Historian'ın eski verileri otomatik olarak S3'e (Ucuz Depolama) taşıması sayesinde bu ay depolama maliyetlerinden %40 tasarruf ettiklerini görür.
3.  **Karar:** Toplantıda bu tasarrufu yeni bir yatırım bütçesi olarak ekibine sunar.

### Journey 6: Sistem Yöneticisi Can - The Silent Guardian (Sessiz Muhafız)

**Persona:** Can (29), Sistem Yöneticisi (IT/OT).
**Hedef:** Sistemin 7/24 ayakta kalması ve siber güvenlik.

**Senaryo:** Fabrika ağında kısa süreli bir kesinti yaşanır.

**Historian ile:**
1.  **Kesinti Anı:** Ağ kesilse bile, sahadaki "Ingestor" servislerinin (Ring Buffer) veriyi tamponlamaya başladığını izleme ekranından görür. Veri kaybı yoktur.
2.  **Otomatik İyileşme:** Ağ geri geldiğinde, Ingestor'ların biriken veriyi NATS JetStream'e hızla boşalttığını ve sistemin senkronize olduğunu izler.
3.  **Güvenlik:** Bir kullanıcının yetkisiz erişim denemesini Audit loglarından tespit eder ve hesabını kilitler.

### Journey Requirements Summary

Bu yolculuklar şu yetenek gereksinimlerini ortaya çıkarmaktadır:

*   **Kestirimci Analiz:** Fiziksel modelleme (Digital Twin) ile anomali tespiti ve "Kalan Ömür" tahmini.
*   **Akıllı Bildirimler:** Sadece eşik aşımı (Threshold) değil, model sapması (Deviation) bazlı uyarılar.
*   **Gelişmiş Görselleştirme:** Gerçek vs. İdeal (Model) verilerinin aynı grafikte (uPlot) karşılaştırılması.
*   **Yüksek Performanslı Sorgu:** Büyük tarihsel verilerin (1 yıl+) çok hızlı (<100ms) görselleştirilmesi.
*   **Denetim İzi (Audit Trail):** Değiştirilemez, imzalı loglama ve kullanıcı aksiyonlarının takibi.
*   **Mobil Erişim:** Bakım ekipleri ve yöneticiler için tablet/mobil uyumlu arayüzler.
*   **Mobil Erişim:** Bakım ekipleri ve yöneticiler için tablet/mobil uyumlu arayüzler.
*   **Dayanıklılık (Resilience):** Ağ kesintilerinde veri tamponlama (Store and Forward) ve otomatik iyileşme.

## Domain-Specific Requirements

### Industrial Automation (OT) Compliance & Regulatory Overview

Historian, yüksek düzeyde regüle edilen endüstrilerde (İlaç, Gıda, Enerji) kritik görevlerde kullanılacağı için, endüstriyel standartlara uyum bir "özellik" değil, varoluşsal bir zorunluluktur. Ürün, kutudan çıktığı haliyle FDA denetimlerine ve ISA alarm standartlarına hazır olmalıdır.

### Key Domain Concerns

*   **Veri Bütünlüğü (Data Integrity):** Kaydedilen verinin değiştirilmediğinin ve silinmediğinin garantisi.
*   **İzlenebilirlik (Traceability):** Kimin, ne zaman, neyi değiştirdiğinin kesin kaydı.
*   **Operasyonel Güvenlik (Safety):** Alarm sisteminin operatörü doğru yönlendirmesi ve yormaması.
*   **Siber Güvenlik (Cybersecurity):** Kritik altyapıların yetkisiz erişimden korunması.

### Compliance Requirements

#### FDA 21 CFR Part 11 (Electronic Records & Signatures)
*   **Elektronik İmza (Re-Authentication):** Kritik set değeri veya alarm limiti değişikliklerinde, sistem kullanıcıdan şifresini tekrar girmesini istemelidir. Bu işlem "Islak İmza" yerine geçer ve Audit Trail'e işlenir.
*   **Oturum Yönetimi (Auto-Logout):** Operatör konsolları, 15 dakika (yapılandırılabilir) hareketsizlik durumunda oturumu otomatik kapatmalıdır.
*   **Şifre Politikası:** Güçlü şifre zorunluluğu, düzenli değişim (Expiry) ve eski şifrelerin tekrar kullanımının engellenmesi (History) zorunludur.

#### ISA 18.2 (Alarm Management)
*   **Full State Machine:** Alarm durumları eksiksiz uygulanmalıdır: `Unack/Active` -> `Ack/Active` -> `Unack/RTN` -> `Ack/RTN` (Latched Alarms desteği ile).
*   **Alarm Shelving (Raflandırma):** Operatörler, arızalı sensörlerden gelen "gürültülü" alarmları belirli bir süre (örn. 4 saat) veya bakım yapılana kadar geçici olarak bastırabilmelidir (Shelve). Bu özellik, alarm yorgunluğunu önlemek için kritiktir.

### Industry Standards & Best Practices

#### IEC 62443 (Industrial Cybersecurity)
*   **Defense in Depth:** Yazılım mimarisi, derinlemesine savunma prensibine göre tasarlanmalıdır.
*   **Zones and Conduits:** Dağıtım dokümantasyonu, Ingestor'ın DMZ'de (Bölge), Engine'in Güvenli İç Ağda (Bölge) olduğu ve aralarındaki NATS trafiğinin (Kanal) şifrelendiği referans mimariyi sunmalıdır.

### Implementation Considerations

*   **NATS Security:** Tüm veri trafiği varsayılan olarak **TLS 1.2+** ile şifrelenmelidir. Servisler arası kimlik doğrulama için **NATS NKEYs** veya Token mekanizması kullanılmalıdır.
*   **Audit Trail Storage:** Denetim kayıtları, değiştirilemez (immutable) bir yapıda saklanmalı ve zincirleme hash (Chained Hash) ile bütünlükleri korunmalıdır.

## Innovation & Novel Patterns

### Detected Innovation Areas

#### 1. Embedded & Democratized Digital Twin (Gömülü ve Demokratikleştirilmiş Dijital İkiz)
*   **Paradigma Değişimi:** Geleneksel Historian'lar "Dikiz Aynası" (Sadece geçmişi gösterir) iken, Historian bir "Kristal Küre"dir.
*   **Çekirdek Özellik:** Rakiplerin pahalı "Add-on" olarak sattığı simülasyon yeteneği, Historian'ın kalbinde (Core) yer alır. Veri, ETL süreçlerine gerek kalmadan, Ingestion anında simülasyon motoruna beslenir.
*   **Slogan:** *"Don't just record history, predict the future."*

#### 2. Continuous Backtesting (Sürekli Geriye Dönük Test)
*   **Güven Mekanizması:** Modelin doğruluğunu kanıtlamak için sistem arka planda sürekli kendini test eder. (Örn: "Son 24 saatte %92 doğruluk").
*   **Şeffaflık:** Kullanıcıya sadece tahmin değil, o tahminin "Güven Skoru" da sunulur.

### Market Context & Competitive Landscape

*   **Mevcut Durum:** AVEVA ve GE gibi devler, Historian ve Simülasyon ürünlerini ayrı ayrı satar ve entegrasyonu zordur.
*   **Mavi Okyanus:** Historian, bu iki dünyayı tek bir mikroservis mimarisinde birleştirerek, KOBİ'lerden dev tesislere kadar herkesin erişebileceği (Demokratik) bir çözüm sunar.

### Validation Approach

*   **Shadow Mode (Gölge Modu):** Dijital İkiz, başlangıçta sadece izleme modunda çalışır. Tahminler gerçek verilerle kıyaslanır ancak operasyona müdahale etmez.
*   **Model Accuracy Score:** Kullanıcı arayüzünde her tahminin yanında, modelin geçmiş performansına dayalı bir doğruluk yüzdesi gösterilir.

### Risk Mitigation

*   **Advisory Only (Sadece Tavsiye):** Dijital İkiz tahminleri asla güvenlik kilitlerini (ESD) veya kritik kontrol döngülerini (PID) doğrudan yönetmez. Sadece operatöre karar desteği sağlar.
*   **Visual Distinction (Görsel Ayrım):** Tahmini veriler asla gerçek veri gibi (düz çizgi) gösterilmez. Her zaman kesikli çizgi veya gölgeli "Güven Aralığı" (Confidence Interval) ile sunulur.
*   **Visual Distinction (Görsel Ayrım):** Tahmini veriler asla gerçek veri gibi (düz çizgi) gösterilmez. Her zaman kesikli çizgi veya gölgeli "Güven Aralığı" (Confidence Interval) ile sunulur.
*   **Disclaimer:** Arayüzde "Tahmini Veri - Operasyonel Karar İçin Doğrulama Gerektirir" uyarısı her zaman görünür.

## Industrial IoT & SaaS Platform Requirements

### Project-Type Overview

Historian, hibrit bir yapıya sahiptir: Uç noktalarda (Edge) çalışan yüksek performanslı bir IoT veri toplayıcısı ve merkezde (Cloud/On-Prem) çalışan çok kiracılı (Multi-tenant) bir SaaS yönetim platformudur. Bu yapı, hem donanım verimliliğini hem de ölçeklenebilirliği zorunlu kılar.

### Technical Architecture Considerations

#### Hardware & Connectivity (Edge Requirements)
*   **Target Hardware:** Linux tabanlı Endüstriyel Gateway'ler (Örn: Siemens IOT2050, Moxa).
*   **Minimum Specs:** 2 Core CPU, 4GB RAM (Rust/NATS verimliliği sayesinde Java rakiplerinin 1/4'ü kaynakla çalışır).
*   **Store and Forward Strategy:** "Hybrid Spill-to-Disk".
    *   Veri önce RAM (Ring Buffer) üzerine yazılır.
    *   Bağlantı kesilirse ve RAM dolarsa, veriler otomatik olarak diske (SSD/SD Kart) şifreli bir kalıcı kuyruk (Persistent Queue) olarak taşar.
    *   **Retention:** Disk dolana kadar saklama (FIFO). Gorilla sıkıştırmasıyla GB'larca alanda haftalarca veri tutulabilir.

#### Multi-Tenancy & Isolation
*   **Logical Multi-tenancy:** Veritabanı ve NATS konuları (Subjects) `TenantID` ile mantıksal olarak izole edilir (Örn: `tenantA.factory1.sensorX`).
*   **Deployment Flexibility:** Bu yapı sayesinde sistem hem "Shared Cloud SaaS" hem de "Dedicated On-Premise" olarak kurulabilir.

### Implementation Considerations

#### API & Integration Strategy
*   **Primary Data API (GraphQL):** Zaman serisi verilerinin esnek sorgulanması için GraphQL kullanılır. İstemciler sadece ihtiyaç duydukları alanları (Timestamp, Value) çekerek ağ yükünü azaltır.
*   **Management API (REST):** Tag yönetimi, kullanıcı işlemleri ve konfigürasyon için standart REST API kullanılır.
*   **Legacy Support:** Eski ERP/MES sistemleri için CSV/Excel dışa aktarma (Export) endpoint'leri sunulur.
*   **Direct DB Access:** Veritabanı bütünlüğünü korumak için doğrudan SQL/JDBC erişimi **yasaktır**.

#### Advanced RBAC Matrix
Standart rollere ek olarak:
*   **Auditor (Denetçi):** "Read-Only" erişim. Sadece Audit Logları ve Raporları görür. Canlı veriyi veya prosesi göremez. (FDA odaklı).
*   **Auditor (Denetçi):** "Read-Only" erişim. Sadece Audit Logları ve Raporları görür. Canlı veriyi veya prosesi göremez. (FDA odaklı).
*   **Service Account (Bot):** Etkileşimsiz (Non-interactive) API kullanıcısı. Sadece belirli endpoint'lere erişebilir, UI girişi yapamaz.

## Project Scoping & Phased Development

### MVP Strategy & Philosophy

**MVP Approach:** Platform MVP
**Philosophy:** "Az ama Öz ve Kaya Gibi Sağlam." Endüstriyel pazarda yarım ürün kabul görmez. Temel fonksiyonlar, uyumluluk ve performans eksiksiz olmalı, ancak yan özellikler (Raporlama, Eğitim) sonraki fazlara bırakılmalıdır.

### MVP Feature Set (Phase 1)

**Core User Journeys Supported:**
*   Operatör Mehmet (İzleme & Alarm)
*   Mühendis Ayşe (Trend Analizi & Excel Export)
*   Kaliteci Zeynep (Audit Trail)
*   Bakımcı Burak (Kestirimci Uyarı - Lite)

**Must-Have Capabilities:**
*   **Data Ingestion:** Modbus TCP, OPC UA, S7 (Siemens) protokolleri.
*   **Storage Engine:** Rust tabanlı LSM-Tree, Gorilla Sıkıştırma, Tiered Storage (Disk dolana kadar).
*   **Visualization:** Gerçek zamanlı Trendler, Dashboard, Alarm Listesi.
*   **Digital Twin (Lite):** Tek değişkenli anomali tespiti. **Sadece Inference (Çalıştırma).** Model eğitimi platform dışındadır.
*   **Compliance:** ISA 18.2 Alarm State Machine (Shelving dahil), FDA Part 11 Audit Trail (Re-auth dahil).
*   **Export:** Yüksek performanslı CSV/Excel dışa aktarma.

### Post-MVP Features

**Phase 2 (Growth):**
*   **Advanced Reporting:** Platform içi sürükle-bırak rapor oluşturucu.
*   **Model Training UI:** Basit modellerin (Lineer Regresyon vb.) platform içinde eğitilmesi.
*   **Live Integration:** ERP/MES sistemleri için REST/GraphQL yazma API'leri.

**Phase 3 (Expansion):**
*   **Federation:** Çoklu fabrika yönetimi (Cloud Dashboard).
*   **AI Marketplace:** 3. parti model geliştiricilerin modellerini satabileceği pazar yeri.
*   **Mobile App:** Native iOS/Android uygulamaları (Push Notification ile).

### Risk Mitigation Strategy

**Technical Risks:**
*   **Risk:** Dijital İkiz'in donanımı yorması.
*   **Mitigation:** "Training Offline, Inference Online" stratejisi ile işlemci yükü minimize edilir.

**Market Risks:**
*   **Risk:** Raporlama eksikliğinin satışa engel olması.
*   **Mitigation:** Mühendislerin Excel alışkanlığına güvenerek, "Piyasadaki en hızlı Excel Export" özelliğini öne çıkarmak.

**Resource Risks:**
*   **Risk:** Rust geliştirici bulma zorluğu.
*   **Mitigation:** Çekirdek ekip (Core) Rust yazarken, UI ve yan servisler (Alarm, Audit) için daha yaygın dillerin (Go, TS/React) kullanılması.

## Functional Requirements

### Data Ingestion & Management
*   **FR-ING-01 (Protocol Support):** Sistem, Modbus TCP, OPC UA ve Siemens S7 protokolleri üzerinden endüstriyel cihazlardan veri toplayabilmelidir.
*   **FR-ING-02 (Store & Forward):** Sistem, merkezi sunucu ile bağlantı koptuğunda verileri önce RAM'de (Ring Buffer), RAM dolarsa yerel diskte (Spill-to-Disk) şifreli olarak saklamalı ve bağlantı geldiğinde otomatik senkronize etmelidir.
*   **FR-ING-03 (Calculated Tags):** Sistem, ham veriler üzerinde matematiksel ve mantıksal işlemler yaparak "Sanal Tag"ler oluşturabilmeli ve bunları gerçek zamanlı olarak kaydedebilmelidir.

### Visualization & Analysis
*   **FR-VIS-01 (Trend Analysis):** Kullanıcılar, seçilen zaman aralığındaki (1 saat - 10 yıl) verileri interaktif grafikler (Zoom/Pan) üzerinde görüntüleyebilmelidir.
*   **FR-VIS-02 (High-Perf Export):** Kullanıcılar, görüntülenen veri setini (maks. 1M satır) .csv veya .xlsx formatında 5 saniyenin altında indirebilmelidir.
*   **FR-VIS-03 (Confidence Interval):** Dijital İkiz tahmin verileri, gerçek verilerden görsel olarak ayrılmalı (kesikli çizgi) ve güven aralığı (gölgeli alan) ile sunulmalıdır.

### Alarm Management (ISA 18.2)
*   **FR-ALM-01 (Alarm Lifecycle):** Sistem, ISA 18.2 durum makinesini (Unack/Active, Ack/Active, Unack/RTN, Ack/RTN) eksiksiz desteklemelidir.
*   **FR-ALM-02 (Alarm Shelving):** Operatörler, belirli bir alarmı geçici bir süre veya koşul sağlanana kadar (Shelve) susturabilmelidir.

### Digital Twin & Prediction
*   **FR-DT-01 (Inference Engine):** Sistem, önceden eğitilmiş model parametrelerini kullanarak gerçek zamanlı veriler üzerinden anomali tespiti ve değer tahmini yapabilmelidir.
*   **FR-DT-02 (Model Accuracy):** Sistem, her tahmin için geçmiş performansa dayalı bir "Güven Skoru" hesaplamalı ve kullanıcıya göstermelidir.

### Audit & Compliance (FDA Part 11)
*   **FR-AUD-01 (Immutable Logging):** Sistem, tüm kritik kullanıcı işlemlerini (Set değeri değişimi, Alarm Shelving vb.) değiştirilemez bir "Audit Trail" günlüğüne kaydetmelidir.
*   **FR-AUD-02 (Re-Authentication):** Kritik işlemlerde sistem, kullanıcının şifresini tekrar girmesini (Elektronik İmza) zorunlu kılmalıdır.

### System Administration & Security
*   **FR-SYS-01 (Tiered Storage):** Sistem, disk doluluk oranı veya veri yaşına göre eski verileri otomatik olarak Soğuk Depolama alanına (S3/MinIO) taşımalıdır.
*   **FR-SYS-01 (Tiered Storage):** Sistem, disk doluluk oranı veya veri yaşına göre eski verileri otomatik olarak Soğuk Depolama alanına (S3/MinIO) taşımalıdır.
*   **FR-SYS-02 (RBAC):** Sistem; Operatör, Mühendis, Yönetici, Denetçi ve Servis Hesabı gibi farklı rollere dayalı erişim kontrolü sağlamalıdır.

## Non-Functional Requirements

### Performance
*   **NFR-PERF-01 (Write Throughput):** Tek bir Ingestion düğümü, saniyede >50.000 olay/saniye yazma hızını kararlı bir şekilde sürdürebilmelidir.
*   **NFR-PERF-02 (Query Latency):** 1 yıllık tarihsel veri sorgusu (Downsampling ile), 100ms'nin altında tamamlanmalıdır.
*   **NFR-PERF-03 (End-to-End Latency):** Verinin sensörden çıkıp ekranda görünmesi arasındaki süre:
    *   **On-Premise:** < 500ms
    *   **Cloud/SaaS:** < 2 saniye
*   **NFR-PERF-04 (Browser Rendering):** Viz arayüzü, 100.000 veri noktasını çizerken 60 FPS akıcılığını korumalı ve tarayıcıyı dondurmamalıdır.

### Reliability & Availability
*   **NFR-REL-01 (Uptime):** Sistem, yıllık %99.99 çalışma süresi (Uptime) sağlamalıdır.
*   **NFR-REL-02 (Zero Data Loss):** Ağ kesintisi durumlarında "Store & Forward" mekanizması sayesinde %0 veri kaybı garanti edilmelidir.
*   **NFR-REL-03 (Data Integrity):** Disk arızası veya ani güç kesintisi durumunda, WAL (Write Ahead Log) sayesinde son 1 saniye hariç veri bütünlüğü korunmalıdır.
*   **NFR-REL-04 (MTTR):** Kritik servisler (Ingestor, Engine), çökme durumunda < 5 saniye içinde otomatik olarak yeniden başlamalıdır.

### Efficiency (Resource Usage)
*   **NFR-EFF-01 (Edge Footprint):** Ingestor servisi, tam yük altında (>50k/s) çalışırken < 500MB RAM ve < %40 CPU (2 Core cihazda) tüketmelidir.

### Security
*   **NFR-SEC-01 (Encryption):** Tüm veri trafiği (Transit) TLS 1.2+ ile, diskteki veriler (At Rest) AES-256 ile şifrelenmelidir.
*   **NFR-SEC-02 (Immutable Audit):** Denetim kayıtları, kriptografik olarak imzalanmalı ve değiştirilemez olmalıdır.

### Scalability
*   **NFR-SCL-01 (Vertical Scaling):** Sistem, eklenen CPU ve RAM kaynağı ile doğru orantılı (lineer) performans artışı göstermelidir.
*   **NFR-SCL-02 (Horizontal Scaling):** NATS kümeleme (Clustering) ile yük birden fazla sunucuya dağıtılabilmelidir.
