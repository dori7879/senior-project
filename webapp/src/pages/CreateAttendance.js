import Footer from '../components/Footer'
import Header from '../components/Header'
import AttendanceForm from '../components/AttendanceForm'
import { useTranslation } from 'react-i18next';

const CreateAttendance = () => {
  const { t } = useTranslation(['translation', 'attendance']);
  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center pb-4'>
          <div className='pt-4 text-xl font-bold text-purple-900'>
            {t('attendance:caption', 'Creating an Attendance Assignment')}
          </div>
        </div>
        <AttendanceForm />
        <small className='mt-4'>{t('attendance:required', '* - Required field')} </small>
      </div>
      <Footer />
    </div>
  )
}

export default CreateAttendance
