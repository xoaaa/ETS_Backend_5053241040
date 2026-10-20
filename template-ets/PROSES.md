# PROSES

## 1. Persona dan masalahnya
(2-3 kalimat dengan kata-katamu sendiri: siapa pengguna persona kamu, masalah apa yang diselesaikan fitur khasnya.)
Pengguna persona A adalah petugas gedung. Beliau sering ditanya mengenai ruangan-ruangan yang kosong ada di jam berapa, sehingga beliau harus membuka daftar peminjaman satu per satu. Fitur ketersediaan slot membantu beliau mengetahui jumlah slot yang tersisa dan slot mana saja yang masih tersedia di tanggal tertentu.

## 2. Rancangan struktur data
(Tulis struct dan slice/map apa yang kamu simpan, dan kenapa dipilih begitu. Rancang ini sebelum menulis kode, lalu catat kalau berubah di tengah jalan.)
Saya menggunakan struct `SumberDaya` dengan field `id`, `nama`, `jenis`, dan `stok`. `stok` adalah jumlah unit yang dimiliki. Struct `Peminjaman` menyimpan `id`, `sumber_daya_id`, `peminjam`, `tanggal`, `slot_mulai`, `slot_selesai`, `jumlah`, dan `status`.
Saya menggunakan slice untuk menyimpan daftar sumber daya dan peminjaman. ID sumber daya dan peminjaman memiliki penghitung masing-masing agar ID tidak digunakan ulang setelah data dihapus. Data delapan slot tidak disimpan, tetapi dihitung setiap kali endpoint ketersediaan dipanggil. Sisa setiap slot dihitung dari stok dikurangi jumlah unit yang sedang dipinjam pada sumber daya, tanggal, dan slot tersebut. Data bersama dilindungi menggunakan mutex.

## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
