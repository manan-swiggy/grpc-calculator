package helper

func Sieve(num int) []int {
	boolArr := make([]bool, num+1)
	for i := 0; i < num+1; i++ {
		boolArr[i] = false
	}

	for i := 2; i*i <= num; i++ {
		if !boolArr[i] {
			for j := i * 2; j <= num; j += i {
				boolArr[j] = true
			}
		}
	}

	primes := []int{}
	for i := 2; i <= num; i++ {
		if !boolArr[i] {
			primes = append(primes, i)
		}
	}

	return primes
}