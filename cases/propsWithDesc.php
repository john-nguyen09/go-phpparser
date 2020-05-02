<?php /**
		* Information about a course that is cached in the course table 'modinfo' field (and then in
		* memory) in order to reduce the need for other database queries.
		*
		* This includes information about the course-modules and the sections on the course. It can also
		* include dynamic data that has been updated for the current user.
		*
		* Use {@link get_fast_modinfo()} to retrieve the instance of the object for particular course
		* and particular user.
		*
		* @property-read int $courseid Course ID
		* @property-read int $userid User ID
		* @property-read array $sections Array from section number (e.g. 0) to array of course-module IDs in that
		*     section; this only includes sections that contain at least one course-module
		* @property-read cm_info[] $cms Array from course-module instance to cm_info object within this course, in
		*     order of appearance
		* @property-read cm_info[][] $instances Array from string (modname) => int (instance id) => cm_info object
		* @property-read array $groups Groups that the current user belongs to. Calculated on the first request.
		*     Is an array of grouping id => array of group id => group id. Includes grouping id 0 for 'all groups'
		*/