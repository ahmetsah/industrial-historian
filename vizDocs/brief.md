1. Vizyon ve Temel Felsefe
"Endüstriyel Verinin Hız Motoru" Tıpkı McMaster-Carr'ın fiziksel parça arama sürecindeki sürtünmeyi (friction) sıfıra indirmesi gibi, bu sistem de milyarlarca satırlık zaman serisi verisine erişimdeki sürtünmeyi sıfıra indirecek.

Motto: "Yükleme ekranı yok. Beklemek yok. Sadece veri."

Estetik Anlayışı: "Boring Technology". Animasyonlar, gölgeler veya büyük boşluklar yok. Yoğun, bilgi odaklı, yüksek kontrastlı ve okunabilir bir arayüz.

2. Hedef Kitle (Personas)
Bakım Mühendisi (Maintenance Eng.): Bir motor arızalandığında, son 24 saatin titreşim verisini 2 saniye içinde görmek ister. Renkli dashboard'lar değil, ham veri ve limit aşımı arar.

Proses Mühendisi: 5 yıllık sıcaklık trendini analiz etmek için CSV indirmek veya Excel ile boğuşmak istemez. Tarayıcıda anında görmek ister.

Entegratör/Developer: API'nin temiz, dokümante edilmiş ve Rust hızında yanıt vermesini bekler.

3. Kritik Ürün Özellikleri (The McMaster Way)
A. Parametrik "Tag" Arama Motoru (Faceted Search)
Standart bir "arama çubuğu" yerine, McMaster'ın sol menüsü gibi dinamik bir filtreleme sistemi.

Taksonomi: Fabrika > Üretim Hattı > Makine > Komponent > Sensör Tipi

Özellik: Kullanıcı "Sıcaklık" filtresini seçtiğinde, sistem anında sadece sıcaklık sensörlerini listeler. Sayfa yenilenmez.

B. "Sıfır Gecikme" Görselleştirme (Viz Layer)
Teknoloji: viz/ klasöründeki yapı, WebGL veya Canvas tabanlı olmalı.

Deneyim: Bir sensör listesinde mouse ile gezinirken (hover), yan tarafta o sensörün son 1 saatlik "kıvılcım grafiği" (sparkline) anında belirmeli (Predictive Prefetching). Tıklayınca detay açılmalı.

C. İçerik Derinliği (Metadata Density)
Bir sensör (Tag) listelendiğinde sadece ismi yazmamalı. Tıpkı McMaster'ın vida detayları gibi:

Mühendislik Birimi (Bar, °C)

Tarama Sıklığı (100ms, 1s)

Veri Tipi (Float32, Bool)

Bağlı olduğu PLC/Cihaz ID'si

Hepsi tek satırda, taranabilir formatta.

4. Teknik Kısıtlar ve Gereksinimler (Repo Analizine Göre)
Core: Rust (crates/historian-core) ile yazılan motor, diskten/bellekten veriyi "stream" ederek okumalı. Pagination yerine "Infinite Scroll" veya "Virtual Scrolling" kullanılmalı (DOM şişkinliğini önlemek için).

API: Go servisleri (go-services), ön yüzün (Frontend) ihtiyaç duyduğu "faset sayılarını" (örn: 'Basınç sensörü: 45 adet') önceden hesaplamalı (Aggregated Counts).

5. Başarı Kriterleri (KPIs)
TTI (Time to Interactive): < 300ms

Sorgu Hızı: 1 milyon veri noktasını çizdirmek < 1 saniye.

Kullanıcı Hızı: Bir mühendisin siteye girip, spesifik bir sensörün geçen haftaki verisini bulup grafiğini açması < 10 saniye.