package pegparser

type ParseResult[T any] struct {
	success bool
	val     T
}

type Parser[T any] func([]byte, int) (ParseResult[T], int)

func fail[T any]() ParseResult[T] {
	var zeroVal T
	return ParseResult[T]{success: false, val: zeroVal}
}

func succeed[T any](val T) ParseResult[T] {
	return ParseResult[T]{success: true, val: val}
}

func Char(c byte) Parser[byte] {
	return func(source []byte, loc int) (ParseResult[byte], int) {
		if loc >= len(source) {
			return fail[byte](), loc
		}
		if c == source[loc] {
			return succeed(c), loc + 1
		}
		return fail[byte](), loc
	}
}

func Digit() Parser[int] {
	return func(source []byte, loc int) (ParseResult[int], int) {
		if loc >= len(source) {
			return fail[int](), loc
		}
		c := source[loc]
		if c >= '0' && c <= '9' {
			return succeed(int(c - '0')), loc + 1
		}
		return fail[int](), loc
	}
}

func Apply[T1 any, T2 any](f func(T1) T2, p Parser[T1]) Parser[T2] {
	return func(source []byte, loc int) (ParseResult[T2], int) {
		res1, loc2 := p(source, loc)
		var res2 ParseResult[T2]
		if res1.success {
			res2.val = f(res1.val)
			res2.success = true
			return res2, loc2
		} else {
			return res2, loc
		}
	}
}

func Seq[T any](ps []Parser[T]) Parser[[]T] {
	return func(source []byte, loc int) (ParseResult[[]T], int) {
		if loc >= len(source) {
			return fail[[]T](), loc
		}
		loc2 := loc
		var res ParseResult[[]T]
		for _, p := range ps {
			var r ParseResult[T]
			r, loc2 = p(source, loc2)
			if !r.success {
				return fail[[]T](), loc
			}
			res.val = append(res.val, r.val)
		}
		res.success = true
		return res, loc2
	}
}

func Choice[T any](ps []Parser[T]) Parser[T] {
	return func(source []byte, loc int) (ParseResult[T], int) {
		if loc >= len(source) {
			return fail[T](), loc
		}
		for _, p := range ps {
			res, loc2 := p(source, loc)
			if res.success {
				return res, loc2
			}
		}
		return fail[T](), loc
	}
}

func Many[T any](p Parser[T]) Parser[[]T] {
	return func(source []byte, loc int) (ParseResult[[]T], int) {
		var res ParseResult[[]T]
		pr, loc := p(source, loc)
		for pr.success {
			res.val = append(res.val, pr.val)
			pr, loc = p(source, loc)
		}
		res.success = true
		return res, loc
	}
}

func Many1[T any](p Parser[T]) Parser[[]T] {
	return func(source []byte, loc int) (ParseResult[[]T], int) {
		pr, loc2 := p(source, loc)
		if !pr.success {
			return fail[[]T](), loc
		}
		var res ParseResult[[]T]
		for pr.success {
			res.val = append(res.val, pr.val)
			pr, loc2 = p(source, loc2)
		}
		res.success = true
		return res, loc2
	}
}

func Maybe[T any](p Parser[T]) Parser[*T] {
	return func(source []byte, loc int) (ParseResult[*T], int) {
		pr, loc2 := p(source, loc)
		var res ParseResult[*T]
		res.success = true
		if pr.success {
			res.val = &pr.val
			return res, loc2
		} else {
			res.val = nil
			return res, loc
		}
	}
}

func And[T1 any, T2 any](p1 Parser[T1], p2 Parser[T2]) Parser[T1] {
	return func(source []byte, loc int) (ParseResult[T1], int) {
		pr1, loc1 := p1(source, loc)
		if !pr1.success {
			return fail[T1](), loc
		}
		pr2, _ := p2(source, loc1)
		if !pr2.success {
			return fail[T1](), loc
		}
		return pr1, loc1
	}
}

func Not[T any](p Parser[T]) Parser[T] {
	return func(source []byte, loc int) (ParseResult[T], int) {
		var res ParseResult[T]
		pr, _ := p(source, loc)
		res.success = !pr.success
		return res, 0
	}
}
