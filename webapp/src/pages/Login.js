import Footer from '../components/Footer'
import Header from '../components/Header'
import LoginForm from '../components/LoginForm'
import { useTranslation } from 'react-i18next';

const Login = () => {
    const { t } = useTranslation(['translation', 'login']);
  return (
    <div>
      <Header />
      <div className='flex items-center justify-center min-h-screen px-4 py-8 bg-purple-100 sm:px-6 lg:px-8'>
        <div className='border border-purple-300'>
          <div className='w-64 max-w-md mx-4'>
            <div>
              <h2 className='mt-6 text-2xl font-extrabold leading-7 text-center text-purple-900 uppercase'>
                {t('login:caption', 'LOGIN')}
              </h2>
            </div>
            <LoginForm />
          </div>
        </div>
      </div>
      <Footer />
    </div>
  )
}

export default Login
