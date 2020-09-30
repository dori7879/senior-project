import React from 'react';

const HomePage = () => {
    return(
        <div className="min-h-screen flex items-start justify-between bg-purple-100 py-8 px-4 sm:px-6 lg:px-8"> 
           <div>
               <img src="https://img.pngio.com/online-education-png-education-png-456_413.png" alt="" className=""/>
           </div>
           <div>
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
    )
}

export default HomePage;