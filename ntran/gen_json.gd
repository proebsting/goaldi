#  gen_json.gd -- create json output from intermediate representation.

procedure json_File(irgen, flagList) {
	local p
	local flag
	local s

	flag := nil
	s := "[\n"
	while p := @irgen do {
		if \flag then s ||:= ",\n"
		flag := "true"
		s ||:= json(p, "")
	}
	s ||:= "\n]"
	return s
}

procedure json_list(p, indent) {
	local s
	local flag
	local i

	s := "["
	flag := nil
	every i := !p do {
		if \flag then {
			s ||:= ","
		}
		flag := "true"
		s ||:= "\n" || indent || "\t" || json(i, indent || "\t")
	}
	s ||:= "\n" || indent || "]"
	return s
}

procedure json_record(p, indent) {
	local s
	local i

	s := "{\n" || indent || "\t\"tag\" : " || image(type(p).name())
	every i := 1 to *p do {
		s ||:= ",\n" || indent || "\t"
		s ||:= image(p.type()[i])
		s ||:= " : "
		s ||:= json(p[i], indent || "\t")
	}
	s ||:= "\n" || indent || "}"
	return s
}

procedure json(p, indent) {
	case type(p) of {
		niltype:		return "null"
		number:			return image(image(p))	# all digits, quoted
		string:			return json_image(string(p))
		list:			return json_list(p, indent)
		set:			return json_list(p, indent)
		ir_Label:		return image(p.value)
		ir_Tmp:			return image(p.name)
		ir_TmpLabel:	return image(p.name)
		ir_TmpClosure:	return image(p.name)
		default:		return json_record(p, indent)
	}
}

procedure json_image(s) {
	/static mapping := table() {
		"\x00"	: `\u0000`,
		"\x01"	: `\u0001`,
		"\x02"	: `\u0002`,
		"\x03"	: `\u0003`,
		"\x04"	: `\u0004`,
		"\x05"	: `\u0005`,
		"\x06"	: `\u0006`,
		"\x07"	: `\u0007`,
		"\b"	: `\b`,
		"\t"	: `\t`,
		"\n"	: `\n`,
		"\v"	: `\u000b`,
		"\f"	: `\f`,
		"\r"	: `\r`,
		"\x0e"	: `\u000e`,
		"\x0f"	: `\u000f`,
		"\x10"	: `\u0010`,
		"\x11"	: `\u0011`,
		"\x12"	: `\u0012`,
		"\x13"	: `\u0013`,
		"\x14"	: `\u0014`,
		"\x15"	: `\u0015`,
		"\x16"	: `\u0016`,
		"\x17"	: `\u0017`,
		"\x18"	: `\u0018`,
		"\x19"	: `\u0019`,
		"\x1a"	: `\u001a`,
		"\e"	: `\u001b`,
		"\x1c"	: `\u001c`,
		"\x1d"	: `\u001d`,
		"\x1e"	: `\u001e`,
		"\x1f"	: `\u001f`,
		`"`		: `\"`,
		`\`		: `\\`,
		"\d"	: `\u007f`,
	}
	local t := `"`
	every local c := !s do
		t ||:= \mapping[c] | c
	return t || `"`
}
