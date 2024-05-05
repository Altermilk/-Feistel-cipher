package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
  rand.Seed(time.Now().UnixNano()) 

  fmt.Print("Введите строку для шифрования: ")
  input := bufio.NewReader(os.Stdin)
  str, _ := input.ReadBytes('\n') // считываем строку для шифрования
  str = []byte(strings.TrimSpace(string(str))) //вырезаем лишние пробелы если они есть

  rounds := 32
  g := generateKeys(8, rounds) // генерируем массив ключей размером 64бит (8 байт)

  encrypted := encrypt(str, g, rounds) // шифруем строку

  fmt.Println("Шифруем : ")
  fmt.Println(encrypted)
  fmt.Println(str)
  decpypted := decrypt([]byte(encrypted), g, rounds) // расшифровываем шифротекст

  fmt.Println("Расшифровываем :", decpypted)

  fmt.Println("Совпадает ли размер входного блока и выходного? :", len(str)==len(decpypted))

  /*fmt.Println("раундовые ключи: ")
  for i := range g{
    fmt.Println(string(g[i]))
  }*/
}

func encrypt(info []byte, k [][]byte, rounds int) string{ // функция шифрования, принимает в себя текст, массив ключей и количество раундов

  n := len(info)/2 // длина правой и левой частей
  balanced := true

  if(len(info)%2 != 0){ // если сеть несбалансированная, то увеличиваем размер частей на 1 и в конец правой добавляем пробел
    n+= len(info)%2;
    balanced = false;
  }

  L := make([]byte, n)
  R := make([]byte, n)

  copy(L, info[:n])
  copy(R, info[n:])
  if(!balanced){
    R[n-1] = byte(' ')
  }

  for i := 0; i<rounds; i++{
    L, R = R, xor(gamma(R, k[i]), L) // тело раунда 
  }

  encrypted := make([]byte, n*2) // соединяем части обратно, снова меняя их местами
    copy(encrypted[:n], R)
    copy(encrypted[n:], L)

    return string(encrypted)
}

func decrypt(info []byte, k [][]byte, rounds int) string{ // функция дешифрования
    n := len(info) / 2
    L := make([]byte, n)
    R := make([]byte, n)
    copy(L, info[:n]) // копируем части в обратном порядке
    copy(R, info[n:])
    for i := 0; i < rounds; i++ {
        R, L = xor(gamma(R, k[rounds-i-1]), L), R // раунд аналогичен шифрованию, но ключи в обратном порядке
    }
    decrypted := make([]byte, n*2)
    copy(decrypted[:n], R) // снова меняем части местами и соединяем 
    copy(decrypted[n:], L)

  return strings.TrimSpace(string(decrypted))
}

func xor(a, b []byte) []byte{ 

  result := make([]byte, len(a))

  for i := range a{
    result[i] = a[i]^b[i]
  }

  return result
}

func gamma(info, key [] byte) [] byte { // функция наложения гаммы
  k := make([]byte, len(info))
  j := 0
  for i := 0; i < len(info); i++ { // создаем гамму той же длины и записываем в нее ключ
    k[i] = key[j]
    if(j == len(key) - 1){  // если текст длиннее ключа, то дублируем гамму до его конца
      j = 0
    }else{
      j++
    }
  }

  return xor(info, k) 
}

func generateKeys(size, rounds int) [][]byte {
  
  rand.Seed(time.Now().UnixNano())
  keys := make([][]byte, rounds)
  for i := 0; i < rounds; {
    keys[i] = make([]byte, size)
    
    for j := 0; j < size; j++ {
      keys[i][j] = byte(rand.Intn(255))
    }

    if (len(keys[i])!=0){
      i++
    }
  }

  return keys
}