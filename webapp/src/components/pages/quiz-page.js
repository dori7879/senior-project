import Footer from '../footer';
import Header from '../header';
import QuizForm from '../quiz-form';
import React from 'react';

const QuizPage = () => {
    return(
        <div>
            <Header />
            <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                <div className="flex justify-center flex-col items-center pb-4">
                    <div className="text-purple-900 font-bold text-xl pt-4">Creating a Quiz</div>
                </div>
                <QuizForm />
                <small className="mt-4">* - Required field </small>  
            </div>
            <Footer />
        </div>
    )
}

export default QuizPage;