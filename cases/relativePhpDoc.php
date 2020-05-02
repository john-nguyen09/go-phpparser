<?php
class TestClass1 {
	/**
	 * @return static
	 */
	public function method1() {}

	/**
	 * @return TestClass2|$this
	 */
	public function method2() {}
}