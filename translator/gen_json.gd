#  gen_json.gd -- create json output from intermediate representation.

procedure json_File(f, irgen) {
	local sep
	f.write("[")
	while ^p := @irgen do {
		f.writes(\sep)
		json(f, p, "")
		sep := ","
	}
	f.write("\n]")
}

procedure json(f, p, indent) {		# write p to f
	case type(p) of {
		niltype:		f.writes("null")
		number:			f.writes(image(image(p)))	# all digits, quoted
		string:			f.writes(json_image(string(p)))
		ir_Label:		f.writes(image(p.value))
		ir_Tmp:			f.writes(image(p.name))
		ir_TmpLabel:	f.writes(image(p.name))
		ir_TmpClosure:	f.writes(image(p.name))
		set:			json_list(f, p, indent)
		list:			json_list(f, p, indent)
		default:		return json_record(f, p, indent)
	}
}

procedure json_list(f, p, indent) {
	local sep
	f.writes("[")
	every ^i := !p do {
		f.writes(\sep | "", "\n", indent, "\t")
		json(f, i, indent || "\t")
		sep := ","
	}
	f.writes("\n", indent, "]")
}

procedure json_record(f, p, indent) {
	f.writes("{\n", indent, "\t\"tag\" : ", image(type(p).name()))
	every ^i := 1 to *p do {
		f.writes(",\n", indent, "\t", image(p.type()[i]), " : ")
		json(f, p[i], indent || "\t")
	}
	f.writes("\n", indent, "}")
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
