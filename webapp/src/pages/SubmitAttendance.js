import { Redirect, useParams } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next';
import { useForm } from "react-hook-form";
import Footer from '../components/Footer'
import Header from '../components/Header'
import AttendanceService from '../services/attendance';

const SubmitAttendance = () => {
  const { t } = useTranslation(['translation', 'attendance']);
  const { randomStr } = useParams()
  const isLoggedIn = useSelector((state) => state.auth.isLoggedIn)
  const { register, handleSubmit } = useForm();
  
  // Attendance related state
  const [courseTitle, setCourseTitle] = useState('')
  const [title, setTitle] = useState('')
  const [mode, setMode] = useState('all')
  // eslint-disable-next-line no-unused-vars
  const [closeDate, setCloseDate] = useState(new Date())
  // eslint-disable-next-line no-unused-vars
  const [attendanceID, setAttendanceID] = useState(-1)

  // Form related state
  const [successful, setSuccessful] = useState(false)

  useEffect(() => {
    AttendanceService.fetchStudentAttendance(randomStr)
      .then((response) => {
        if (response.data) {
          setAttendanceID(response.data.ID)
          setMode(response.data.Mode)
          setCourseTitle(response.data.CourseTitle)
          setTitle(response.data.Title)
          setCloseDate(response.data.ClosedAt)
        }
      })
      .catch((error) => {
        alert(error.message)
      })
  }, [randomStr, register]);

  const onSubmit = (data) =>  {
    AttendanceService.submitAttendance({...data, AttendanceID: parseInt(attendanceID)})
      .then(
        (response) => {
          setSuccessful(true);
        },
        (error) => {
          alert(error.message);
        }
      )
  };
  
  if (mode === 'registered' && !isLoggedIn) {
    return <Redirect to='/login' />
  }
  
  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center w-3/4 pb-4'>
          <div className='pt-2 text-xl font-bold text-purple-900'>{t('attendance:attendance', 'Attendance')}</div>
          <div className='flex flex-row items-center w-full border border-purple-300 rounded-t items '>
            <div className='w-full p-4 bg-purple-300 border border-purple-300 rounded-tl'>
              <h2 className='block px-4 pt-1 mb-3 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                <strong>{t('attendance:coursetitle', 'Course Title')}:</strong>{' '}
                <span className='text-purple-900'>{courseTitle}</span>
              </h2>
              <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                <strong>{t('attendance:attTitle', 'Attendance Title:')}</strong>{' '}
                <span className='text-purple-900'>{title}</span>
              </h2>
            </div>
            <div>
              {isLoggedIn ? null : (
                <div className='flex flex-row items'>
                  <label className='block px-4 pt-1 mb-4 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                    {t('attendance:fullname', 'Full Name*:')}
                  </label>
                  <input
                    name="StudentFullName"
                    ref={register}
                    className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                    type='text'
                    placeholder={t('attendance:entername','Enter your full name')}
                  />
                </div>
              )}
              <div className='flex flex-row items-center items'>
                <h2 className='block px-4 pt-1 mb-2 text-xs font-bold tracking-wide text-gray-700 uppercase'>
                  <strong>{t('attendance:time', 'Time Remaining:')}</strong>{' '}
                  <span className='text-purple-900'>15{t('attendance:min', 'min')}</span>
                </h2>
              </div>
            </div>
          </div>
          <div className='w-full p-4 bg-purple-300 border border-purple-300 rounded-b'>
            {successful ? (
              <div className='form-group'>
                <div className='w-full pt-2 text-xl font-bold text-center text-purple-900 border border-purple-300 rounded'>
                {t('attendance:submitted', 'Submitted')}
                </div>
              </div>
            ) : (
  
              <form onSubmit={handleSubmit(onSubmit)} >
                <div className='flex flex-col pb-2 mx-4'>
                  <h1 className='block px-4 pt-1 mt-4 mb-2 text-xs font-bold tracking-wide text-center text-purple-900 uppercase'>
                  {t('attendance:pin', 'PIN')}
                  </h1>
                  <input
                    name="PIN"
                    ref={register}
                    className='px-1 mb-3 mr-4 text-xs leading-tight text-gray-700 border border-purple-400 rounded focus:outline-none focus:bg-white'
                    type='text'
                    placeholder={t('attendance:enterpin','Enter PIN')}
                  />
                </div>
                <div className='flex justify-center'>
                  <button
                    type='submit'
                    className='relative flex justify-center w-full px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
                  >
                    {t('attendance:submitatt', 'Submit')}
                  </button>
                </div>
              </form>
            )}
          </div>
        </div>
        <small className=''>{t('attendance:required', '* - Required field')} </small>
      </div>
      <Footer />
    </div>
  )
}

export default SubmitAttendance