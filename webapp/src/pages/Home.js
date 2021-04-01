import Footer from '../components/Footer'
import Header from '../components/Header'
import { Link } from 'react-router-dom'
import Steps from '../components/Steps'
import { useTranslation } from 'react-i18next';

const Home = () => {
  const { t } = useTranslation(['translation', 'home']);
  return (
    <div>
      <Header />
      <div className='flex flex-col justify-center min-h-screen px-4 py-2 bg-purple-100 items sm:px-6 lg:px-8 '>
        <div className='flex items-start justify-center pt-6'>
          <div className='flex justify-around pb-12'>
            <div className='mr-10 '>
              <img
                src='https://img.pngio.com/online-education-png-education-png-456_413.png'
                alt=''
                className=''
              />
            </div>
            <div className='ml-10 lg:mt-20 sm:mt-10'>
              <Link to='/homework'>
                <button
                  type='button'
                  className='relative flex items-center justify-center w-full h-10 px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
                >
                  {t('home:hw', 'Create homework page')}
                </button>
              </Link>
              <Link to='/quiz'>
                <button
                  type='button'
                  className='relative flex items-center justify-center w-full h-10 px-2 py-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
                >
                  {t('home:quiz', 'Create quiz page')}
                </button>
              </Link>
              <Link to='/attendance'>
                <button
                  type='button'
                  className='relative flex items-center justify-center w-full h-10 px-2 pt-1 mb-2 text-sm font-medium leading-4 text-purple-200 transition duration-150 ease-in-out bg-purple-800 border border-transparent rounded-md hover:bg-purple-500 focus:outline-none'
                >
                  {t('home:attendance', 'Create attendance page')}
                </button>
              </Link>
            </div>
          </div>
        </div>
        <div className='flex flex-col items-center justify-center pb-8 border-t-2 border-purple-200'>
          <div className='pt-4 pb-2 text-xl font-bold text-purple-900'>
            {t('home:caption', ' Introducing EasySubmit')}
          </div>
          <div className='w-2/3 bg-purple-200 border-2 border-purple-400 rounded-lg'>
            <p className='flex justify-center mx-4 my-2 text-justify text-gray-700'>
              {t('home:intro', 'A progressive web application (PWA) that allows teachers/instructors to post homework (quiz or attendance) simply by opening the website and generating a link for the new homework page. This link could then be shared with students to submit their solutions. Furthermore, the proposed tool allows creating different types of quizzes with an auto grading system. Another functionality of the system is in-class attendance checking mechanism using QR codes or a generated link ID. In the second case, the teacher/instructor can tick one or more required functionalities and this check will be added into the link.')}{' '}
            </p>
          </div>
        </div>
        <Steps />
      </div>
      <Footer />
    </div>
  )
}

export default Home
