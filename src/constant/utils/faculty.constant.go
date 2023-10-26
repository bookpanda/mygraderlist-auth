package utils

type Faculty struct {
	FacultyEN string `json:"faculty_en"`
	FacultyTH string `json:"faculty_th"`
}

var Faculties = map[string]Faculty{
	"20": {
		FacultyEN: "Graduate School",
		FacultyTH: "บัณฑิตวิทยาลัย",
	},
	"21": {
		FacultyEN: "Faculty of Engineering",
		FacultyTH: "คณะวิศวกรรมศาสตร์",
	},
	"22": {
		FacultyEN: "Faculty of Arts",
		FacultyTH: "คณะอักษรศาสตร์",
	},
	"23": {
		FacultyEN: "Faculty of Science",
		FacultyTH: "คณะวิทยาศาสตร์",
	},
	"24": {
		FacultyEN: "Faculty of Political Science",
		FacultyTH: "คณะรัฐศาสตร์",
	},
	"25": {
		FacultyEN: "Faculty of Architecture",
		FacultyTH: "คณะสถาปัตยกรรมศาสตร์",
	},
	"26": {
		FacultyEN: "Faculty of Commerce And Accountancy",
		FacultyTH: "คณะพาณิชยศาสตร์และการบัญชี",
	},
	"27": {
		FacultyEN: "Faculty of Education",
		FacultyTH: "คณะครุศาสตร์",
	},
	"28": {
		FacultyEN: "Faculty of Communication Arts",
		FacultyTH: "คณะนิเทศศาสตร์",
	},
	"29": {
		FacultyEN: "Faculty of Economics",
		FacultyTH: "คณะเศรษฐศาสตร์",
	},
	"30": {
		FacultyEN: "Faculty of Medicine",
		FacultyTH: "คณะแพทยศาสตร์",
	},
	"31": {
		FacultyEN: "Faculty of Veterinary Science",
		FacultyTH: "คณะสัตวแพทยศาสตร์",
	},
	"32": {
		FacultyEN: "Faculty of Dentistry",
		FacultyTH: "คณะทันตแพทยศาสตร์",
	},
	"33": {
		FacultyEN: "Faculty of Pharmaceutical Sciences",
		FacultyTH: "คณะเภสัชศาสตร์",
	},
	"34": {
		FacultyEN: "Faculty of Law",
		FacultyTH: "คณะนิติศาสตร์",
	},
	"35": {
		FacultyEN: "Faculty of Fine And Applied Arts",
		FacultyTH: "คณะศิลปกรรมศาสตร์",
	},
	"36": {
		FacultyEN: "Faculty of Nursing",
		FacultyTH: "คณะพยาบาลศาสตร์",
	},
	"37": {
		FacultyEN: "Faculty of Allied Health Sciences",
		FacultyTH: "คณะสหเวชศาสตร์",
	},
	"38": {
		FacultyEN: "Faculty of Psychology",
		FacultyTH: "คณะจิตวิทยา",
	},
	"39": {
		FacultyEN: "Faculty of Sports Science",
		FacultyTH: "คณะวิทยาศาสตร์การกีฬา",
	},
	"40": {
		FacultyEN: "School of Agricultural Resources",
		FacultyTH: "วิทยาลัยประชากรศาสตร์",
	},
	"51": {
		FacultyEN: "College of Population Studies",
		FacultyTH: "วิทยาลัยประชากรศาสตร์",
	},
	"53": {
		FacultyEN: "College of Public Health Sciences",
		FacultyTH: "วิทยาลัยวิทยาศาสตร์สาธารณสุข",
	},
	"55": {
		FacultyEN: "Language Institute",
		FacultyTH: "สถาบันภาษา",
	},
	"56": {
		FacultyEN: "School of Integrated Innovation",
		FacultyTH: "สถาบันนวัตกรรมบูรณาการ",
	},
	"58": {
		FacultyEN: "Sasin Graduate Institute of Business Administion",
		FacultyTH: "สถาบันบัณฑิตบริหารธุรกิจ ศศินทร์ฯ",
	},
	"99": {
		FacultyEN: "Other University",
		FacultyTH: "มหาวิทยาลัยอื่น",
	},
	"01": {
		FacultyEN: "The Sirindhorn Thai Language Institute",
		FacultyTH: "สถาบันภาษาไทยสิรินธร",
	},
	"02": {
		FacultyEN: "Office of Academic Affairs",
		FacultyTH: "ศูนย์การศึกษาทั่วไป",
	},
}
