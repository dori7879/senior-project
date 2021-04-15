import Footer from '../components/Footer'
import Header from '../components/Header'
import QuizForm from '../components/QuizForm'

const Quiz = () => {
  return (
    <div>
      <Header />
      <div className='flex flex-col items-center min-h-screen px-4 py-2 bg-purple-100 justify-top sm:px-6 lg:px-8 '>
        <div className='flex flex-col items-center justify-center pb-4'>
          <div className='pt-4 text-xl font-bold text-purple-900'>
            Creating a Quiz
          </div>
        </div>
        <QuizForm />
        <small className='mt-4'>* - Required field </small>
      </div>
      <Footer />
    </div>
  )
}

export default Quiz
