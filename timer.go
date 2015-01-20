// 2014-12-26 Adam Bryt

// Program timer (countdown timer) odlicza czas do zera. Co sekundę
// wyświetla na terminalu czas pozostały do końca, a na końcu
// generuje sygnał dźwiękowy.
//
// Sposób użycia:
//
//	timer [duration]
//
// gdzie duration jest czasem w formacie time.Duration, np: "2h30m15s";
// domyślnie 25m.
//
// Uwagi:
//
// Terminal musi obsługiwać sekwencje kontrolne ANSI, które są
// używane do wyświetlania czasu w lewym górnym rogu terminala.
// Sygnał dźwiękowy jest generowany przez wysłanie do terminala
// znaku '\a'. Jeśli terminal ma ustawioną opcję popOnBell to po
// wygenerowaniu dźwięku okno terminala jest umieszczane na wierzchu.
//
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	defaultTime = 25 * time.Minute // domyślny czas gdy podano opcji -t
	stepTime    = 1 * time.Second  // co ile wyświetlać pozostały czas

	bell  = "\a"
	clear = "\x1b[2J" // clear screen
	home  = "\x1b[H"  // move cursor to upper left corner

	beepNum   = 5                      // ilość beepnięć na końcu
	beepSpace = 300 * time.Millisecond // odstęp między beepnięciami
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: timer [duration]\n")
	fmt.Fprintf(os.Stderr, "\tduration jest czasem w formacie '1h20m15s'\n")
	flag.PrintDefaults()
	os.Exit(2)
}

// print wyświetla na terminalu czas pozostały do końca, w formacie
// time.Duration. Terminal musi obsługiwać sekwencje sterujące ANSI.
func print(t time.Duration) {
	fmt.Print(clear)
	fmt.Print(home)
	fmt.Println(t)
}

// beep generuje sygnał dźwiękowy przez wysłanie na terminal sekwencji
// znaków '\a'.
func beep() {
	for i := 0; i < beepNum; i++ {
		fmt.Print(bell)
		time.Sleep(beepSpace)
	}
}

// timer odlicza czas od t do 0 drukując co sekundę czas pozostały
// do końca. Po upływie czasu t generuje sygnał dźwiękowy.
func timer(t time.Duration) {
	if t < stepTime {
		print(t)
		time.Sleep(t)
		beep()
		return
	}

	ticker := time.NewTicker(stepTime)
	defer ticker.Stop()
	for {
		print(t)
		<-ticker.C
		t -= stepTime
		if t < stepTime {
			print(t)
			time.Sleep(t)
			beep()
			return
		}
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var t time.Duration = defaultTime
	if flag.NArg() > 0 {
		var err error
		t, err = time.ParseDuration(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "timer: %s\n", err)
			usage()
		}
	}
	timer(t)
}
