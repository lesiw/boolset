package a

func use(...any) {}

func clean(keys []string) {
	seen := make(map[string]bool) // want `seen is a set and can be map\[T\]struct\{\}`
	for _, k := range keys {
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = true
	}
	use(len(seen))
}

func bareRead(keys []string) {
	seen := make(map[string]bool) // want `seen is a set and can be map\[T\]struct\{\}`
	for _, k := range keys {
		if seen[k] {
			continue
		}
		seen[k] = true
	}
}

func negRead(keys []string) {
	seen := make(map[string]bool) // want `seen is a set and can be map\[T\]struct\{\}`
	for _, k := range keys {
		if !seen[k] {
			seen[k] = true
		}
	}
}

func falseWrite() {
	m := make(map[string]bool)
	m["a"] = true
	m["b"] = false
	if m["a"] {
		use(1)
	}
}

func deleteAllowed() {
	m := make(map[string]bool) // want `m is a set and can be map\[T\]struct\{\}`
	m["a"] = true
	if m["a"] {
		delete(m, "a")
	}
	clear(m)
}

func threeState() {
	m := make(map[string]bool)
	m["a"] = true
	v, ok := m["a"]
	use(v, ok)
}

func dynamic(v bool) {
	m := make(map[string]bool)
	m["a"] = v
	use(m)
}

func escapes() map[string]bool {
	m := make(map[string]bool)
	m["a"] = true
	return m
}

func valuesRead() {
	m := make(map[string]bool)
	m["a"] = true
	for k, v := range m {
		use(k, v)
	}
}

func elsewhere() {
	m := make(map[string]bool) // want `m is a set and can be map\[T\]struct\{\}`
	m["a"] = true
	use(m["a"])
}

func param(did map[string]bool) {
	if !did["a"] {
		did["a"] = true
	}
}

type holder struct {
	seen map[string]bool
}

func (h *holder) mark() {
	if !h.seen["a"] {
		h.seen["a"] = true
	}
}

func literal() {
	m := map[string]bool{"a": true, "b": true} // want `m is a set and can be map\[T\]struct\{\}`
	if !m["a"] {
		m["b"] = true
	}
}

func literalFalse() {
	m := map[string]bool{"a": false}
	if !m["a"] {
		m["b"] = true
	}
}

func reassigned() {
	m := map[string]bool{"a": true} // want `m is a set and can be map\[T\]struct\{\}`
	if !m["a"] {
		m["b"] = true
	}
	m = make(map[string]bool, 4)
	if !m["c"] {
		m["c"] = true
	}
}

func aliasedSource(src map[string]bool) {
	m := src
	if !m["a"] {
		m["b"] = true
	}
}

type stringSet map[string]bool

func outflow(p *stringSet) {
	m := make(map[string]bool)
	if !m["a"] {
		m["b"] = true
	}
	*p = m
}

var pkgSet = map[string]bool{"a": true}

func packageScope() {
	if !pkgSet["a"] {
		pkgSet["b"] = true
	}
}

func shadowedTrue() {
	true := false
	m := map[string]bool{}
	m["a"] = true
	_ = true
}

func ranged(ms []map[string]bool) {
	for _, m := range ms {
		if !m["a"] {
			m["b"] = true
		}
	}
}

func varCommaOK() {
	m := map[string]bool{"a": true}
	var v, ok = m["a"]
	_, _ = v, ok
}

func varBlankCommaOK() {
	m := map[string]bool{"a": true} // want `m is a set and can be map\[T\]struct\{\}`
	var _, ok = m["a"]
	_ = ok
}

func shadowedMake(src map[string]bool) {
	make := func(int) map[string]bool { return src }
	m := make(0)
	if !m["a"] {
		m["b"] = true
	}
}

type flag bool

func namedElem() {
	m := make(map[string]flag) // want `m is a set and can be map\[T\]struct\{\}`
	if !m["a"] {
		m["b"] = true
	}
}

func namedMap() {
	m := make(stringSet) // want `m is a set and can be map\[T\]struct\{\}`
	if !m["a"] {
		m["b"] = true
	}
}

func escapesOnly() {
	m := map[string]bool{"a": true}
	use(m)
}

func shadowedLiteral() {
	true := false
	m := map[string]bool{"a": true}
	if !m["b"] {
		use(m["a"])
	}
	_ = true
}

func aliasedOut() {
	m := make(map[string]bool)
	var x = m
	if !m["a"] {
		x["b"] = true
	}
}

func tuple() (int, bool) { return 0, true }

func tupleStore() {
	m := map[string]bool{}
	var x int
	x, m["a"] = tuple()
	use(x)
}

func supply() (m map[string]bool) {
	m = make(map[string]bool)
	if !m["a"] {
		m["b"] = true
	}
	return m
}

func keysOnly() {
	m := map[string]bool{"a": true} // want `m is a set and can be map\[T\]struct\{\}`
	for k := range m {
		use(k)
	}
}
