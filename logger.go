package main

import "log"

/*
	 func prNormal(fmtString string, args ...interface{}) {
		log.Printf(fmtString, args...)
	}
*/
func prInfo(fmtString string, args ...interface{}) {
	fs := "*** INFO: " + fmtString
	log.Printf(mkGreen(fs), args...)
}

/*
	 func prWarn(fmtString string, args ...interface{}) {
		fs := "*** WARN: " + fmtString
		log.Printf(mkYellow(fs), args...)
	}

	func prError(fmtString string, args ...interface{}) {
		log.Printf(mkRed(fmtString), args...)
	}

	func prVerbose(fmtString string, args ...interface{}) {
		if vp.GetBool("Verbose") {
			log.Printf(fmtString, args...)
		}
	}
*/
func prVerboseInfo(fmtString string, args ...interface{}) {
	if vp.GetBool("Verbose") {
		prInfo(fmtString, args...)
	}
}

func prDebug(fmtString string, args ...interface{}) {
	if vp.GetBool("Debug") {
		fs := "*** DEB: " + fmtString
		log.Printf(mkRed(fs), args...)
	}
}
