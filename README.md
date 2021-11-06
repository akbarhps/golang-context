# Golang Context

Sumber Tutorial:
[Udemy](https://www.udemy.com/course/pemrograman-go-lang-pemula-sampai-mahir/learn/lecture/24915656#questions) |
[Slide](https://docs.google.com/presentation/d/1WhJvRpKPWq7LY9P6fMN93vkjKa1bJwBQebbieKdefPw/edit)

## Pengenalan Context
---

- Context merupakan sebuah data yang membawa value, sinyal cancel, sinyal timeout dan sinyal deadline
- Context biasanya dibuat per request (misal setiap ada request masuk ke server web melalui http request)
- Context digunakan untuk mempermudah kita meneruskan value, dan sinyal antar proses


### Kenapa Context Perlu Dipelajari?

- Context di Golang biasa digunakan untuk mengirim data request atau sinyal ke proses lain
- Dengan menggunakan context, ketika kita ingin membatalkan semua proses, kita cukup mengirim sinyal ke context, maka secara otomatis semua proses akan dibatalkan
- Hampir semua bagian di Golang memanfaatkan context, seperti database, http server, http client, dan lain-lain
- Bahkan di Google sendiri, ketika menggunakan Golang, context wajib digunakan dan selalu dikirim ke setiap function yang dikirim


### Cara Kerja Context

![Cara Kerja Context](https://user-images.githubusercontent.com/69947442/140592884-d7f648a3-88ea-4364-a7e8-0b51c2db4ef1.png)


### Package Context

- Context direpresentasikan di dalam sebuah interface Context
- interface Context terdapat dalam package context
- https://golang.org/pkg/context/ 


### Interface Context

![Interface Context](https://user-images.githubusercontent.com/69947442/140592922-390e6f49-181e-4ff0-834b-a30cf56e709e.png)


### Membuat Context

- Karena Context adalah sebuah interface, untuk membuat context kita butuh sebuah struct yang sesuai dengan kontrak interface Context
- Namun kita tidak perlu membuatnya secara manual
- Di Golang package context terdapat function yang bisa kita gunakan untuk membuat Context


### Function Membuat Context

| Function               | Keterangan                                                                                                                                                                                           |
| ---------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `context.Background()` | Membuat context kosong. Tidak pernah dibatalkan, tidak pernah timeout, dan tidak memiliki value apapun. Biasanya digunakan di main function atau dalam test, atau dalam awal proses request terjadi. |
| `context.TODO()`       | Membuat context kosong seperti Background(), namun biasanya menggunakan ini ketika belum jelas context apa yang ingin digunakan                                                                      |


### Kode: Membuat Context

```go
func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background) // context.Background

	todo := context.TODO()
	fmt.Println(todo) // context.TODO
}
```


## Parent dan Child Context
---

- Context menganut konsep parent dan child
- Artinya, saat kita membuat context, kita bisa membuat child context dari context yang sudah ada
- Parent context bisa memiliki banyak child, namun child hanya bisa memiliki satu parent context
- Konsep ini mirip dengan pewarisan di pemrograman berorientasi object


### Diagram Parent dan Child Context

![Diagram Parent dan Child Context](https://user-images.githubusercontent.com/69947442/140593404-5e9ee710-fb60-417e-802a-8e7060af6887.png)


### Hubungan Antara Parent dan Child Context

- Parent dan Child context akan selalu terhubung
- Saat nanti kita melakukan misal pembatalan context A, maka semua child dan sub child dari context A akan ikut dibatalkan
- Namun jika misal kita membatalkan context B, hanya context B dan semua child dan sub child nya yang dibatalkan, parent context B tidak akan ikut dibatalkan
- Begitu juga nanti saat kita menyisipkan data ke dalam context A, semua child dan sub child nya bisa mendapatkan data tersebut
- Namun jika kita menyisipkan data di context B, hanya context B dan semua child dan sub child nya yang mendapat data, parent context B tidak akan mendapat data


### Immutable

- Context merupakan object yang Immutable, artinya setelah Context dibuat, dia tidak bisa diubah lagi
- Ketika kita menambahkan value ke dalam context, atau menambahkan pengaturan timeout dan yang lainnya, secara otomatis akan membentuk child context baru, bukan merubah context tersebut


### Cara Membuat Child Context

- Cara membuat child context ada banyak caranya, yang akan kita bahas di materi-materi selanjutnya


## Context With Value
---

- Pada saat awal membuat context, context tidak memiliki value
- Kita bisa menambah sebuah value dengan data Pair (key - value) ke dalam context
- Saat kita menambah value ke context, secara otomatis akan tercipta child context baru, artinya original context nya tidak akan berubah sama sekali
- Untuk membuat menambahkan value ke context, kita bisa menggunakan function context.WithValue(parent, key, value)


### Kode: Context With Value

```go
func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextC := context.WithValue(contextA, "c", "C")
	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA) // context.Background
	fmt.Println(contextB) // context.Background.WithValue(type string, val B)
	fmt.Println(contextC) // context.Background.WithValue(type string, val C)
	fmt.Println(contextD) // context.Background.WithValue(type string, val B).WithValue(type string, val D)
	fmt.Println(contextE) // context.Background.WithValue(type string, val B).WithValue(type string, val E)
	fmt.Println(contextF) // context.Background.WithValue(type string, val C).WithValue(type string, val F)
}
```

### Kode: Context Get Value

```go
func TestContextGetValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextC := context.WithValue(contextA, "c", "C")
	contextF := context.WithValue(contextC, "f", "F")

	fmt.Println(contextA.Value("b")) // nil
	fmt.Println(contextB.Value("a")) // nil
	fmt.Println(contextC.Value("a")) // nil
	fmt.Println(contextD.Value("b")) // B
	fmt.Println(contextE.Value("b")) // B
	fmt.Println(contextF.Value("b")) // nil
}
```


## Context With Cancel
---

- Selain menambahkan value ke context, kita juga bisa menambahkan sinyal cancel ke context
- Kapan sinyal cancel diperlukan dalam context?
- Biasanya ketika kita butuh menjalankan proses lain, dan kita ingin bisa memberi sinyal cancel ke proses tersebut
- Biasanya proses ini berupa goroutine yang berbeda, sehingga dengan mudah jika kita ingin membatalkan eksekusi goroutine, kita bisa mengirim sinyal cancel ke context nya
- Namun ingat, goroutine yang menggunakan context, tetap harus melakukan pengecekan terhadap context nya, jika tidak, tidak ada gunanya
- Untuk membuat context dengan cancel signal, kita bisa menggunakan function `context.WithCancel(parent)`


### Kode: Contoh Goroutine Leak

```go
func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	// goroutine leak
	// go func() {
	// 	defer close(destination)

	// 	for counter := 1; ; counter++ {
	// 		destination <- counter
	// 	}
	// }()

	// no goroutine leak
	go func() {
		defer close(destination)

		for counter := 1; ; counter++ {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter: ", n)

		// Counter:  1
		// Counter:  2
		// Counter:  3
		// Counter:  4
		// Counter:  5
		// Counter:  6
		// Counter:  7
		// Counter:  8
		// Counter:  9
		// Counter: 10

		if n == 10 {
			break
		}
	}

	cancel() // send cancel signal to context

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
	// When Goroutine Leak: Total Goroutine 3
	// When Goroutine No Leak: Total Goroutine 2
}
```


## Context With Timeout
---

- Selain menambahkan value ke context, dan juga sinyal cancel, kita juga bisa menambahkan sinyal cancel ke context secara otomatis dengan menggunakan pengaturan timeout
- Dengan menggunakan pengaturan timeout, kita tidak perlu melakukan eksekusi cancel secara manual, cancel akan otomatis di eksekusi jika waktu timeout sudah terlewati
- Penggunaan context dengan timeout sangat cocok ketika misal kita melakukan query ke database atau http api, namun ingin menentukan batas maksimal timeout nya
- Untuk membuat context dengan cancel signal secara otomatis menggunakan timeout, kita bisa menggunakan function context.WithTimeout(parent, duration)


### Kode: Context With Timeout

```go
func CreateSlowCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)

		for counter := 1; ; counter++ {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				time.Sleep(time.Second)
			}
		}
	}()

	return destination
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()

	destination := CreateSlowCounter(ctx)
	for n := range destination {
		fmt.Println("Counter: ", n)

		// Counter:  1
		// Counter:  2
		// Counter:  3
		// Counter:  4
		// Counter:  5
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2
}
```


## Context With Deadline
---

- Selain menggunakan timeout untuk melakukan cancel secara otomatis, kita juga bisa menggunakan deadline
- Pengaturan deadline sedikit berbeda dengan timeout, jika timeout kita beri waktu dari sekarang, kalo deadline ditentukan kapan waktu timeout nya, misal jam 12 siang hari ini
- Untuk membuat context dengan cancel signal secara otomatis menggunakan deadline, kita bisa menggunakan function context.WithDeadline(parent, time)


### Kode: Context With Deadline

```go
func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2

	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
	defer cancel()

	destination := CreateSlowCounter(ctx)
	for n := range destination {
		fmt.Println("Counter: ", n)

		// Counter:  1
		// Counter:  2
		// Counter:  3
		// Counter:  4
		// Counter:  5
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine()) // Total Goroutine 2
}
```