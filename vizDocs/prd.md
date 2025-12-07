1. Giriş ve Amaç
Bu proje, endüstriyel zaman serisi verileri için McMaster-Carr felsefesine sahip (aşırı hızlı, faydacı, sürtünmesiz) bir tarihçi (historian) ve görselleştirme motoru geliştirmeyi amaçlar.

Temel Hedef: Bir mühendisin, milyarlarca veri noktası arasından spesifik bir sensörü bulup, son 24 saatlik verisini görselleştirmesi 10 saniyenin altında gerçekleşmelidir.

2. Kullanıcı Personaları
Bakım Mühendisi (Maintenance Engineer - "Fix-it Felix"):

Amacı: Arızalı bir ekipmanın neden durduğunu anlamak.

Davranışı: Spesifik bir "Tag ID" arar. Hız ister. Renkli grafikler değil, limit aşımı (threshold breach) görmek ister.

Proses Mühendisi (Process Engineer - "Data Diana"):

Amacı: Üretim verimliliğini artırmak için uzun vadeli trendleri (1 yıl+) incelemek.

Davranışı: Karşılaştırma yapar. Üst üste birden fazla sensörü bindirir (overlay). CSV indirmekten nefret eder.

Sistem Entegratörü (Integrator - "API Ali"):

Amacı: Bu veriyi başka bir sisteme (ERP/MES) aktarmak.

Davranışı: Swagger/OpenAPI dokümantasyonu arar.

3. Fonksiyonel Gereksinimler
3.1. Arama ve Navigasyon (The "McMaster" Search)
FR-01 Fasetli Arama: Kullanıcılar sensörleri hiyerarşik özelliklere göre filtreleyebilmelidir. (Fabrika > Hat > Makine > Sensör Tipi).

FR-02 Anlık Sonuç (Instant Feedback): Filtre seçimi yapıldığında sayfa yenilenmemeli, sonuç listesi <100ms içinde güncellenmelidir.

FR-03 Metadata Gösterimi: Arama sonuçlarında Tag ID'nin yanında; Birim, Tarama Sıklığı ve Veri Tipi tek satırda gösterilmelidir.

3.2. Veri Görselleştirme (Zero-Latency Viz)
FR-04 Predictive Prefetch (Öngörülü Yükleme): Kullanıcı fareyi (mouse) bir sensör satırının üzerine getirdiğinde (hover), sistem son 1 saatin özet verisini arka planda çekmelidir.

FR-05 Sparklines: Liste görünümünde, her sensörün yanında son duruma dair mini bir grafik (sparkline) yer almalıdır.

FR-06 Derin Zoom: Ana grafikte fare tekerleği ile zoom yapıldığında, veri "downsampling" (örnekleme azaltma) olmadan, Rust backend'den dinamik olarak detaylandırılmalıdır.

3.3. Veri Yönetimi
FR-07 Çoklu Veri Tipi: Float, Integer, Boolean ve String veri tiplerini desteklemelidir.

FR-08 Export: Seçili aralığın CSV/JSON olarak dışa aktarımı.

4. Teknik ve Teknik Olmayan Gereksinimler (Non-Functional)
4.1. Performans
NFR-01: "Time to Interactive" (TTI) < 300ms olmalıdır.

NFR-02: 1 milyon noktalı bir sorgu, görselleştirme dahil < 1 saniye sürmelidir.

4.2. Altyapı ve Dağıtım
NFR-03 (Kritik): Sistem tamamen Dockerize edilmelidir. docker-compose up komutu ile Database, Backend (Rust/Go) ve Frontend ayağa kalkmalıdır.

NFR-04: Backend, Rust (crates/historian-core) tabanlı olmalı, API katmanı Go (go-services) ile servis edilmelidir.

4.3. UX/UI Standartları
NFR-05: "Boring UI". Standart sistem fontları, yüksek kontrast, gereksiz animasyon yok.

NFR-06: Klavye kısayolları (J/K ile listede gezme, / ile arama) desteklenmelidir.