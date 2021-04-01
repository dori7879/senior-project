import Footer from '../components/Footer'
import Header from '../components/Header'
import RegistrationForm from '../components/RegistrationForm'
import { useTranslation } from 'react-i18next';

const Registration = () => {
  const { t } = useTranslation(['translation', 'registration']);
  return (
    <div>
      <Header />
      <div className='flex items-center justify-center min-h-screen px-4 py-8 bg-purple-100 sm:px-6 lg:px-8'>
        <div className='border border-purple-300'>
          <div className='w-64 max-w-md mx-4'>
            
            <RegistrationForm />
          </div>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default Registration
