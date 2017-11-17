package hld

func NewStaffDirCmd() StaffDirCmd {
	v := StaffDirCmd{}
	v.Params = make(map[string]interface{})
	v.Command = "findStaff"
	v.Params["jsData"] = true
	return v
}

//StaffDirCmd Staff Directory command
type StaffDirCmd struct {
	HSJsonCommand
}

//SetStaffCode Set staffcode mask
func (cmd *StaffDirCmd) SetStaffCode(staffCode string) {
	cmd.Params["staffCode"] = staffCode
}

//StaffCode Return staffcode
func (cmd *StaffDirCmd) StaffCode() (string, bool) {
	switch s := cmd.Params["staffCode"].(type) {
	case string:
		return s, true
	default:
		return "", false
	}
}

//SetSurname Set surname mask
func (cmd *StaffDirCmd) SetSurname(surname string) {
	cmd.Params["surname"] = surname
}

//SetLastname Set lastname mask, actually it's the first name
func (cmd *StaffDirCmd) SetLastname(lastname string) {
	cmd.Params["lastname"] = lastname
}

//SetChinese Set chinese mask
func (cmd *StaffDirCmd) SetChinese(chinese string) {
	cmd.Params["chinese"] = chinese
}

//SetShowTerminated show terminated staff
func (cmd *StaffDirCmd) SetShowTerminated(b bool) {
	cmd.Params["showTerminate"] = b
}

//Setname Set name mask
func (cmd *StaffDirCmd) Setname(name string) {
	cmd.Params["name"] = name
}

//SetName Set name mask
func (cmd *StaffDirCmd) SetName(name string) {
	cmd.Params["name"] = name
}

//SetExt Set extension mask
func (cmd *StaffDirCmd) SetExt(ext string) {
	cmd.Params["ext"] = ext
}

//SetLoginID Set login id mask
func (cmd *StaffDirCmd) SetLoginID(loginid string) {
	cmd.Params["loginID"] = loginid
}

//SetDeptCode Set department code mask
func (cmd *StaffDirCmd) SetDeptCode(deptCode string) {
	cmd.Params["deptCode"] = deptCode
}

//StaffDirItem data row for each record
type StaffDirItem struct {
	Sex          string     `json:"sex"`
	Region       string     `json:"_region"`
	EmpID        int        `json:"emp_id"`
	JoinDate     HSJSONTime `json:"joinDate"`
	Chinese      string     `json:"chinese"`
	Country      string     `json:"country"`
	SubDivCode   string     `json:"subDivCode"`
	CompCode     string     `json:"compCode"`
	Status       string     `json:"status"`
	DeptCode     string     `json:"deptCode"`
	Extension    string
	Domain       string
	Surname      string
	Lastname     string
	NewStaffCode string
	LastEmpDate  HSJSONTime
	LoginID      string
	GradeCode    string
	Email        string
	StaffCode    string
	Title        string
	Name         string
	Area         string
	IsVolunteer  bool
	Alias        string
	ProbaEnd     HSJSONTime
}

/*
{
	"sex": "M",
	"_region": "HK",
	"emp_id": 11278,
	"joinDate": 1280678400000,
	"chinese": "陳頴倫",
	"country": "NT",
	"subDivCode": "AD2",
	"compCode": "HRE",
	"status": "A",
	"deptCode": "EDP",
	"extension": "5124",
	"domain": "hq.hk.hld",
	"surname": "Chan",
	"postCode": "323",
	"lastUpdate": 1451988028000,
	"lastname": "Wing Lun",
	"newStaffCode": "HR09903",
	"lastEmpDate": null,
	"loginId": "EDPcwlc",
	"gradeCode": "C15",
	"email": "wl.chan@hld.com",
	"staffCode": "HR9903",
	"title": "MR",
	"loaded": false,
	"name": "Chan Wing Lun, Alan",
	"area": "NW02",
	"isVolunteer": false,
	"alias": "Alan",
	"probaEnd": 1288627200000
}*/
