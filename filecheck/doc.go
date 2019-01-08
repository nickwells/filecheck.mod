/*

Package filecheck provides a means of checking that certain constraints apply
to a file-system object. There are checks on the existence of the file and on
it's status. The caller can opt to apply the checks to the object itself or,
if it is a symbolic link, to the object being linked to. The default
behaviour will check the status of the linked-to object.

*/
package filecheck
