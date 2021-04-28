import { Redirect, useParams } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { useEffect, useState } from 'react'
import moment from "moment";
import Footer from '../components/Footer'
import Header from '../components/Header'
import AttendanceService from '../services/attendance';
import QRCode from 'qrcode.react'

const ListAttSubmissions = () => {
  const { randomStr } = useParams()
  const role = useSelector((state) => state.auth.role)
  const isLoggedIn = useSelector((state) => state.auth.isLoggedIn)

  // Attendance related state
  // eslint-disable-next-line no-unused-vars
  const [attendanceID, setAttendanceID] = useState(-1)
  const [subs, setAttSubs] = useState([])
  const [title, setTitle] = useState('')
  const [pin, setPIN] = useState('')
  const [studentLink, setStudentLink] = useState('')
  // eslint-disable-next-line no-unused-vars
  const [closeDate, setCloseDate] = useState(new Date())
  const [showing, setShowing] = useState(false)
  const [timeout, setTimeoutValue] = useState("7")

  useEffect(() => {
    AttendanceService.fetchTeacherAttendance(randomStr)
      .then((response) => {
        if (response.data) {
          setAttendanceID(response.data.ID)
          setTitle(response.data.Title)
          setStudentLink(response.data.StudentLink)
          setCloseDate(response.data.ClosedAt)
          setPIN(response.data.PIN)
          if ('Submissions' in response.data) {
            setAttSubs(response.data.Submissions)
          }
        }
      })
      .catch((error) => {
        alert(error.message)
      })
  }, [randomStr])

  const handleShowing = () => {
    if (!showing) {
      setShowing(true)
    }

    setTimeout(() => {
      AttendanceService.renewStudentPIN(attendanceID)
        .then((response) => {
          if (response.data) {
            setShowing(false)
            setPIN(response.data.pin)
          }
        })
        .catch((error) => {
          alert(error.message)
        })
    }, +timeout*1000)
  }

  const downloadQR = () => {
    const canvas = document.getElementById("studentLinkQR");
    const pngUrl = canvas
      .toDataURL("image/png")
      .replace("image/png", "image/octet-stream");
    let downloadLink = document.createElement("a");
    downloadLink.href = pngUrl;
    downloadLink.download = "studentLinkQR.png";
    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);
  };

  if (isLoggedIn && role === 'student') {
    return <Redirect to='/' />
  }

  const submittedAtts = subs.length > 0 ? subs.length : 0

  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center pb-4'></div>
        <div className='pt-4 text-xl font-bold text-purple-900'>
          Successful Attendance Submissions
        </div>

        <div className='w-3/4 p-4 bg-purple-300 border border-purple-300 rounded'>
          <div className='flex flex-row'>
            <div className='flex flex-col w-3/4'>
              <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                <strong>Attendance Title:</strong>{' '}
                <span className='text-purple-900'>{title}</span>
              </h2>

              <div className="flex flex-col items-center justify-center">
                {/* BIG STUDENT LINK HASH */}
                <h2 className="text-4xl"><strong>{studentLink}</strong></h2>
                {/* BIG QR CODE */}
                <div className="mt-10">
                  <QRCode
                    id="studentLinkQR"
                    value={"/api/v1/attendances/shared/" + studentLink + "/student"}
                    size={290}
                    level={"H"}
                  />
                  <div className="flex justify-center mt-2">
                    <a className="text-xs" href="#a" onClick={downloadQR}> Download QR </a>
                  </div>
                </div>
                {/* button with modal that shows PIN code */}
                <div className="flex flex-col items-center justify-center mt-10">
                  <div className='flex flex-row items'>
                    <label className='block px-4 pt-1 mb-4 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                      Set PIN Renew Timeout (in seconds)
                    </label>
                    <input
                      name="pin"
                      value={timeout}
                      className='px-1 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                      type='number'
                      onChange={(e) => setTimeoutValue(e.target.value)}
                    />
                  </div>
                  <button className="mt-5" onClick={() => handleShowing()}>Show PIN</button>
                  { showing 
                      ? <h2 className='block px-4 pt-1 mb-2 text-4xl font-bold tracking-wide text-gray-700 uppercase'>
                          <strong>PIN:</strong>{' '}
                          <span className='text-purple-900'>{pin}</span>
                        </h2>
                      : null
                  }
                </div>
              </div>

              <div className='flex flex-col ml-2 items '>
                {subs.map((sub, index) => (
                  <div
                    key={index}
                    className='flex flex-col mb-2 border border-purple-700 rounded'
                  >
                    <div className='flex flex-col pb-2 items'>
                      <p className='block px-4 pt-1 text-xs tracking-wide text-gray-700'>
                        <strong>Student Name:</strong>{' '}
                        <span className='text-purple-900'>
                          {sub.StudentFullName ? sub.StudentFullName :
                            sub.Student != null ? sub.Student.FirstName + " " + sub.Student.LastName :
                            ""
                          }
                        </span>
                      </p>
                      <p className='block px-4 pt-1 mb-2 text-xs tracking-wide text-gray-700'>
                        <strong>Submitted at:</strong>{' '}
                        <span className='text-purple-900'>
                          {moment(sub.SubmittedAt, moment.ISO_8601)}
                        </span>
                      </p>
                      <p className='block px-4 pt-1 mb-2 text-xs tracking-wide text-gray-700'>
                        <strong>{sub.Present ? 'PRESENT' : 'NOT PRESENT'}</strong>
                      </p>
                    </div>
                  </div>
                ))}
              </div>

            </div>

            <div className='flex flex-col '>
              <p className='block px-4 pt-1 text-xs tracking-wide text-gray-700'>
                <strong>
                  <span className='font-bold text-purple-800'>
                    {submittedAtts}
                  </span>{' '}
                  Attendances Submitted
                </strong>{' '}
              </p>
              <p className='block px-4 pt-1 text-xs tracking-wide text-gray-700'>
                <strong></strong>{' '}
              </p>
            </div>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default ListAttSubmissions
