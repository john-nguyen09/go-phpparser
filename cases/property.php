<?php

class ClassWithProperties {
    private static $staticProp1;
    protected static $staticProp2;
    public static $staticProp2 = array();
    static $staticProp3, $staticProp4;
    private static $staticProp5, $staticProp6;

    private $prop1;
    protected $prop2;
    public $prop3;
    var $prop4, $prop5;
    protected $prop6, $prop7;

    /**
     * @var ClassWithProperties
     */
    public $prop5;
}
