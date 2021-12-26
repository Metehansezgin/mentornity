


:4000/api/getData ---> Veritabanındaki kayıtların hepsini getirir.

:4000/api/addData ---> Kayıt ekler ve veritabanındaki verilerle karşılaştırarak scor hesabı yapar ve geri döner.

kayıt eklerken kullanılacak json veri yapısı;
  {
    "Name": "isim",
    "City": "şehir",
    "Category": "kategori"
  }
  
  dönen deger;
  
  [
    {
    "Name": "isim",
    "City": "şehir",
    "Category": "kategori"
    "Score":"skor"
    }
  ]
