package proaction

test_containsRef {
  containsRef("a", ["b"], "a", "b")
}

test_containsRef {
  not containsRef("a", ["b"], "c", "b")
}

test_containsRef {
  not containsRef("a", ["1", "b"], "c", "b")
}

test_containsRef {
  containsRef("a1", ["b"], "a1", "b")
}

test_containsRef {
  not containsRef("a1", ["b"], "a1", "c")
}

test_containsRef {
  not containsRef("a1", ["b"], "c1", "d")
}
