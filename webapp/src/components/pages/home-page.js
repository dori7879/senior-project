import React from 'react';

const HomePage = () => {
    return(
        <div className="min-h-screen flex items flex-col justify-center bg-purple-100 py-2 px-4 sm:px-6 lg:px-8">
            <div className="min-h-screen flex items-start justify-between"> 
                <div className="flex justify-around">
                    <div className=" mr-10">
                        <img src="https://img.pngio.com/online-education-png-education-png-456_413.png" alt="" className=""/>
                    </div>
                    <div className="lg:mt-20 ml-10 sm:mt-10">
                            <button type="button" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-600 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                Create homework page
                            </button>
                            <button type="button" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-600 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                Create quiz page
                            </button>
                            <button type="button" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-600 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                Create attendance page
                            </button>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default HomePage;