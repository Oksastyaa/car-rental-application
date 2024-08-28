package pkg

import (
	"fmt"
	"math/rand"
)

// List of Indonesian names
var indonesianNames = []string{
	"Ahmad", "Ayu", "Budi", "Citra", "Dewi", "Eko", "Fitri", "Gita", "Hadi", "Indah",
	"Joko", "Kartini", "Lestari", "Maya", "Nina", "Omar", "Putri", "Qori", "Rizki", "Sari",
	"Tono", "Ujang", "Vina", "Wahyu", "Yuli", "Zainal", "Dimas", "Siti", "Andi", "Dwi",
	"Bagus", "Farah", "Gilang", "Hendra", "Ira", "Jasmine", "Kiki", "Lia", "Mira", "Nia",
	"Oki", "Rina", "Salsa", "Tia", "Ulfa", "Veri", "Winda", "Yana", "Zulfi",
}

var indonesianSurnames = []string{
	"Pratama", "Wijaya", "Putra", "Sari", "Wulandari", "Kurniawan", "Saputra", "Santoso", "Siregar", "Simanjuntak",
	"Ramadhan", "Wibowo", "Maulana", "Gunawan", "Permadi", "Rahayu", "Rahman", "Wibisono", "Purnomo", "Nugraha",
	"Setiawan", "Utomo", "Haryanto", "Ananda", "Susanto", "Subagyo", "Kusuma", "Fauzi", "Hartono", "Aulia",
}

var indonesianCities = []string{
	"Jakarta",
	"Bandung",
	"Surabaya",
	"Medan",
	"Semarang",
	"Yogyakarta",
	"Denpasar",
	"Makassar",
	"Palembang",
	"Banjarmasin",
	"Balikpapan",
	"Padang",
	"Malang",
	"Manado",
	"Pontianak",
	"Banda Aceh",
	"Bandar Lampung",
	"Ambon",
	"Batam",
	"Batu",
	"Bau-Bau",
	"Bekasi",
	"Bengkulu",
	"Binjai",
	"Bitung",
	"Blitar",
	"Bogor",
	"Bontang",
	"Bukittinggi",
	"Cilegon",
	"Cimahi",
	"Cirebon",
	"Denpasar",
	"Depok",
	"Dumai",
	"Gorontalo",
	"Gunungsitoli",
	"Jakarta Barat",
	"Jakarta Pusat",
	"Jakarta Selatan",
	"Jakarta Timur",
	"Jakarta Utara",
	"Jambi",
	"Jayapura",
	"Kediri",
	"Kendari",
	"Kotamobagu",
	"Kupang",
	"Langsa",
	"Lhokseumawe",
	"Lubuk Linggau",
	"Madiun",
	"Magelang",
	"Makassar",
	"Malang",
	"Manado",
	"Mataram",
	"Medan",
	"Metro",
	"Mojokerto",
	"Padang",
	"Padang Panjang",
	"Padang Sidempuan",
	"Pagar Alam",
	"Palangka Raya",
	"Palembang",
	"Palopo",
	"Palu",
	"Pangkal Pinang",
	"Parepare",
	"Pariaman",
	"Pasuruan",
	"Payakumbuh",
	"Pekalongan",
	"Pekanbaru",
	"Pematang Siantar",
	"Pontianak",
	"Prabumulih",
	"Probolinggo",
	"Sabang",
	"Salatiga",
	"Samarinda",
	"Sawah Lunto",
	"Semarang",
	"Serang",
	"Sibolga",
	"Singkawang",
	"Solok",
	"Sorong",
	"Subulussalam",
	"Sukabumi",
	"Sungaipenuh",
	"Surabaya",
	"Surakarta",
	"Tangerang",
	"Tangerang Selatan",
	"Tanjung Balai",
	"Tanjung Pinang",
	"Tarakan",
	"Tasikmalaya",
	"Tebing Tinggi",
	"Tegal",
	"Ternate",
	"Tidore Kepulauan",
	"Tomohon",
	"Tual",
	"Yogyakarta",
	"Aceh Barat",
	"Aceh Barat Daya",
	"Aceh Besar",
	"Aceh Jaya",
	"Aceh Selatan",
	"Aceh Singkil",
	"Aceh Tamiang",
	"Aceh Tengah",
	"Aceh Tenggara",
	"Aceh Timur",
	"Aceh Utara",
	"Agam",
	"Alor",
	"Asahan",
	"Asmat",
	"Badung",
	"Balangan",
	"Banggai",
	"Banggai Kepulauan",
	"Bangka",
	"Bangka Barat",
	"Bangka Selatan",
	"Bangka Tengah",
	"Bangkalan",
	"Bangli",
	"Banjar",
	"Banjarnegara",
	"Bantaeng",
	"Bantul",
	"Banyuasin",
	"Banyumas",
	"Banyuwangi",
	"Barito Kuala",
	"Barito Selatan",
	"Barito Timur",
	"Barito Utara",
	"Barru",
	"Batang",
	"Batang Hari",
	"Bener Meriah",
	"Bengkalis",
	"Bengkayang",
	"Bengkulu Selatan",
	"Bengkulu Tengah",
	"Bengkulu Utara",
	"Berau",
	"Biak Numfor",
	"Bima",
	"Bireuen",
	"Blora",
	"Boalemo",
	"Bojonegoro",
	"Bolaang Mongondow",
	"Bolaang Mongondow Selatan",
	"Bolaang Mongondow Timur",
	"Bolaang Mongondow Utara",
	"Bombana",
	"Bondowoso",
	"Bone",
	"Bone Bolango",
	"Boven Digoel",
	"Boyolali",
	"Brebes",
	"Buleleng",
	"Bulukumba",
	"Bulungan",
	"Bungo",
	"Buol",
	"Buru",
	"Buru Selatan",
	"Buton",
	"Buton Utara",
	"Ciamis",
	"Cianjur",
	"Cilacap",
	"Deiyai",
	"Deli Serdang",
	"Demak",
	"Dharmasraya",
	"Dogiyai",
	"Dompu",
	"Donggala",
	"Empat Lawang",
	"Ende",
	"Enrekang",
	"Fakfak",
	"Flores Timur",
	"Garut",
	"Gayo Lues",
	"Gianyar",
	"Gorontalo Utara",
	"Gowa",
	"Gresik",
	"Grobogan",
	"Gunung Kidul",
	"Gunung Mas",
	"Halmahera Barat",
	"Halmahera Selatan",
	"Halmahera Tengah",
	"Halmahera Timur",
	"Halmahera Utara",
	"Hulu Sungai Selatan",
	"Hulu Sungai Tengah",
	"Hulu Sungai Utara",
	"Humbang Hasundutan",
	"Indragiri Hilir",
	"Indragiri Hulu",
	"Indramayu",
	"Intan Jaya",
	"Jayawijaya",
	"Jember",
	"Jembrana",
	"Jeneponto",
	"Jepara",
	"Jombang",
	"Kaimana",
	"Kampar",
	"Kapuas",
	"Kapuas Hulu",
	"Karanganyar",
	"Karangasem",
	"Karawang",
	"Karimun",
	"Karo",
	"Katingan",
	"Kaur",
	"Kayong Utara",
	"Kebumen",
	"Kediri",
	"Keerom",
	"Kendal",
	"Kepahiang",
	"Kepulauan Anambas",
	"Kepulauan Aru",
	"Kepulauan Mentawai",
	"Kepulauan Meranti",
	"Kepulauan Sangihe",
	"Kepulauan Seribu",
	"Kepulauan Siau Tagulandang Biaro",
	"Kepulauan Sula",
	"Kepulauan Talaud",
	"Kepulauan Yapen",
	"Kerinci",
	"Ketapang",
	"Klaten",
	"Klungkung",
	"Kolaka",
	"Kolaka Utara",
	"Konawe",
	"Konawe Selatan",
	"Konawe Utara",
	"Kotabaru",
	"Kotawaringin Barat",
	"Kotawaringin Timur",
	"Kuantan Singingi",
	"Kubu Raya",
	"Kudus",
	"Kulon Progo",
	"Kuningan",
	"Kutai Barat",
	"Kutai Kartanegara",
	"Kutai Timur",
	"Labuhan Batu",
	"Labuhan Batu Selatan",
	"Labuhan Batu Utara",
	"Lahat",
	"Lamandau",
	"Lamongan",
	"Lampung Barat",
	"Lampung Selatan",
	"Lampung Tengah",
	"Lampung Timur",
	"Lampung Utara",
	"Landak",
	"Langkat",
	"Lanny Jaya",
	"Lebak",
	"Lebong",
	"Lembata",
	"Lima Puluh Kota",
	"Lingga",
	"Lombok Barat",
	"Lombok Tengah",
	"Lombok Timur",
	"Lombok Utara",
	"Lumajang",
	"Luwu",
	"Luwu Timur",
	"Luwu Utara",
	"Madiun",
	"Magelang",
	"Magetan",
	"Majalengka",
	"Majene",
	"Malinau",
	"Maluku Barat Daya",
	"Maluku Tengah",
	"Maluku Tenggara",
	"Maluku Tenggara Barat",
	"Mamasa",
	"Mamberamo Raya",
	"Mamberamo Tengah",
	"Mamuju",
	"Mamuju Utara",
	"Mandailing Natal",
	"Manggarai",
	"Manggarai Barat",
	"Manggarai Timur",
	"Manokwari",
	"Manokwari Selatan",
	"Mappi",
	"Maros",
	"Maybrat",
	"Melawi",
	"Merangin",
	"Merauke",
	"Mesuji",
	"Mimika",
	"Minahasa",
	"Minahasa Selatan",
	"Minahasa Tenggara",
	"Minahasa Utara",
	"Morowali",
	"Muara Enim",
	"Muaro Jambi",
	"Muko Muko",
	"Muna",
	"Murung Raya",
	"Musi Banyuasin",
	"Musi Rawas",
	"Nabire",
	"Nagan Raya",
	"Nagekeo",
	"Natuna",
	"Nduga",
	"Ngada",
	"Nganjuk",
	"Ngawi",
	"Nias",
	"Nias Barat",
	"Nias Selatan",
	"Nias Utara",
	"Nunukan",
	"Ogan Ilir",
	"Ogan Komering Ilir",
	"Ogan Komering Ulu",
	"Ogan Komering Ulu Selatan",
	"Ogan Komering Ulu Timur",
	"Pacitan",
	"Padang Lawas",
	"Padang Lawas Utara",
	"Padang Pariaman",
	"Pakpak Bharat",
	"Pandeglang",
	"Pangandaran",
	"Pangkajene Kepulauan",
	"Paniai",
	"Parigi Moutong",
	"Pasaman",
	"Pasaman Barat",
	"Paser",
	"Pasuruan",
	"Pati",
	"Pegunungan Arfak",
	"Pegunungan Bintang",
	"Pekalongan",
	"Pelalawan",
	"Pemalang",
	"Penajam Paser Utara",
	"Pesawaran",
	"Pesisir Barat",
	"Pesisir Selatan",
	"Pidie",
	"Pidie Jaya",
	"Pinrang",
	"Pohuwato",
	"Polewali Mandar",
	"Ponorogo",
	"Pontianak",
	"Poso",
	"Pringsewu",
	"Probolinggo",
	"Pulang Pisau",
	"Pulau Morotai",
	"Puncak",
	"Puncak Jaya",
	"Purbalingga",
	"Purwakarta",
	"Purworejo",
	"Raja Ampat",
	"Rejang Lebong",
	"Rembang",
	"Rokan Hilir",
	"Rokan Hulu",
	"Rote Ndao",
	"Sabu Raijua",
	"Sambas",
	"Samosir",
	"Sampang",
	"Sanggau",
	"Sarmi",
	"Sarolangun",
	"Sekadau",
	"Selayar",
	"Seluma",
	"Semarang",
	"Seram Bagian Barat",
	"Seram Bagian Timur",
	"Serang",
	"Serdang Bedagai",
	"Seruyan",
	"Siak",
	"Sidenreng Rappang",
	"Sidoarjo",
	"Sigi",
	"Sijunjung",
	"Sikka",
	"Simalungun",
	"Simeulue",
	"Sinjai",
	"Sintang",
	"Situbondo",
	"Sleman",
	"Solok",
	"Solok Selatan",
	"Soppeng",
	"Sorong",
	"Sorong Selatan",
	"Sragen",
	"Subang",
	"Sukabumi",
	"Sukamara",
	"Sukoharjo",
	"Sumba Barat",
	"Sumba Barat Daya",
	"Sumba Tengah",
	"Sumba Timur",
	"Sumbawa",
	"Sumbawa Barat",
	"Sumedang",
	"Sumenep",
	"Supiori",
	"Sragen",
	"Sumedang",
	"Sumenep",
	"Sungaipenuh",
	"Sintang",
	"Sleman",
	"Surabaya",
	"Surakarta",
	"Tabalong",
	"Tabanan",
	"Takalar",
	"Tambrauw",
	"Tana Tidung",
	"Tana Toraja",
	"Tanah Bumbu",
	"Tanah Datar",
	"Tanah Laut",
	"Tangerang",
	"Tanjung Balai",
	"Tanjung Jabung Barat",
	"Tanjung Jabung Timur",
	"Tanjung Pinang",
	"Tapanuli Selatan",
	"Tapanuli Tengah",
	"Tapanuli Utara",
	"Tapin",
	"Tarakan",
	"Tasikmalaya",
	"Tebing Tinggi",
	"Tegal",
	"Teluk Bintuni",
	"Teluk Wondama",
	"Temanggung",
	"Ternate",
	"Tidore Kepulauan",
	"Timor Tengah Selatan",
	"Timor Tengah Utara",
	"Toba Samosir",
	"Tojo Una-Una",
	"Toli-Toli",
	"Tolikara",
	"Tomohon",
	"Toraja Utara",
	"Trenggalek",
	"Tual",
	"Tuban",
	"Tulang Bawang",
	"Tulang Bawang Barat",
	"Tulungagung",
	"Wajo",
	"Wakatobi",
	"Waropen",
	"Way Kanan",
	"Wonogiri",
	"Wonosobo",
	"Yahukimo",
	"Yalimo",
}

// RandomIndonesianName generates a random Indonesian name
func RandomIndonesianName() string {
	firstName := indonesianNames[rand.Intn(len(indonesianNames))]
	lastName := indonesianSurnames[rand.Intn(len(indonesianSurnames))]
	return fmt.Sprintf("%s %s", firstName, lastName)
}

// RandomCity generates a random Indonesian city
func RandomCity() string {
	return indonesianCities[rand.Intn(len(indonesianCities))]
}

// RandomEmail generates a random email
func RandomEmail() string {
	firstName := indonesianNames[rand.Intn(len(indonesianNames))]
	lastName := indonesianSurnames[rand.Intn(len(indonesianSurnames))]
	separators := []string{"_", ".", "-"}
	separator := separators[rand.Intn(len(separators))]
	return fmt.Sprintf("%s%s%s@gmail.com", firstName, separator, lastName)
}
