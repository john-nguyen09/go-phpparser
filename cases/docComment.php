<?php

/**
 * moodlelib.php - Moodle main library, this is not a version 1.1
 *
 * Main library file of miscellaneous general-purpose Moodle functions.
 * Other main libraries:
 *  - weblib.php      - functions that produce web output
 *  - datalib.php     - functions that access the database
 *
 * @version 1.1.0
 * @author John Citizen <john.citizen@example.com>
 * @author But no author just text but this is recognised as author
 * @package    core
 * @subpackage lib
 * @copyright  1999 onwards Martin Dougiamas  http://dougiamas.com
 * @license    http://www.gnu.org/copyleft/gpl.html GNU GPL v3 or later
 */

/** @var DateTime|int[] $var1 */
/**
 * @method string getString()
 * @method void setInteger(integer $integer)
 * @method setString(integer $integer) But containing the word static results
 *                                     in weird tokens
 * @method static string staticGetter()
 * @method void superMethod(array|DateTime|database|mixed $superParam) This method does
 * a lot of things automagically
 */

 /**
  * @deprecated
  * @deprecated 1.0.0
  * @deprecated No longer used by internal code and not recommended.
  * @deprecated 1.0.0 No longer used by internal code and not recommended.
  */

/**
 * @global database|database2|database[] $DB
 * @link http://example.com/my/bar Documentation of Foo.
 */
global $DB;

/**
 * @global \namespace\database_driver $DRIVER
 */

/**
 * Counts the number of items in the provided array.
 *
 * @param mixed[] $items Array structure to count the elements of.
 *
 * @throws InvalidArgumentException if the provided argument is not of type
 * 'array'.
 * @return int Returns the number of elements.
 */
function count(array $items)
{
}

/**
 * @property string $myProperty
 * @property string
 * @property-read string $myProperty
 * @property-write $myProperty
 */
/** @var string|null Should contain a description */
/**
  * @var string $name        Should contain a description
  * @var string $description Should contain a description
  */
