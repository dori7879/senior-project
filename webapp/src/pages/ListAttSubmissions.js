import { Redirect, useParams } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { useEffect, useState } from 'react'
import Footer from '../components/Footer'
import Header from '../components/Header'
import AttendanceService from '../services/attendance';

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
  // eslint-disable-next-line no-unused-vars
  const [closeDate, setCloseDate] = useState(new Date())

  useEffect(() => {
    AttendanceService.fetchTeacherAttendance(randomStr)
      .then((response) => {
        if (response.data) {
          setAttendanceID(response.data.ID)
          setTitle(response.data.Title)
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
              <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                <strong>PIN:</strong>{' '}
                <span className='text-purple-900'>{pin}</span>
              </h2>

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
                          {sub.SubmittedAt}
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
