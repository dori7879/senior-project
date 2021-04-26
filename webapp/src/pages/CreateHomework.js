import Footer from '../components/Footer'
import Header from '../components/Header'
import HwForm from '../components/HwForm'
import { useTranslation } from 'react-i18next';

const CreateHomework = () => {
  const { t } = useTranslation(['translation', 'homework']);
  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center pb-4'>
          <div className='pt-4 text-xl font-bold text-purple-900'>
            {t('homework:caption', 'Creating a Homework Assignment')}
          </div>
        </div>
        <HwForm />
        <small className='mt-4'>{t('homework:required', '* - Required field')} </small>
      </div>
      <Footer />
    </div>
  )
}

export default CreateHomework
