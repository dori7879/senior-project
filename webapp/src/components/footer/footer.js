import './footer.css';

import React from 'react';

const Footer = () => {
    return(
        <div className="bg-purple-200 ">
            <footer className="text-gray-800 body-font">
                <div className="container px-24 py-12 mx-auto flex md:items-center lg:items-start md:flex-row md:flex-no-wrap flex-wrap flex-col">
                    <div className="w-32 flex-shrink-0 md:mx-0 mx-auto text-center md:text-left md:mt-0 mt-10 ">
                        <a className="flex title-font font-medium items-center md:justify-start justify-center text-gray-900">
                            <img className="w-8 h-8" src="fav.ico" alt="app logo"/>
                            <span className="text-md ml-1 text-purple-900 font-bold">EasySubmit</span>
                        </a>
                        <p className="mt-1  ml-10 text-xs text-gray-600 text-justify">A simple and robust system that makes the interactions as intuitive as possible.</p>
                    </div>
                    <div className="flex-grow flex flex-wrap md:pr-20 -mb-10 md:text-left text-center order-first">
                        <div className="lg:w-1/3 md:w-1/2 w-full px-4">
                            <h2 className="title-font font-medium text-purple-900 tracking-widest text-sm mb-3">ABOUT US</h2>
                            <nav className="list-none mb-10">
                            <li>
                                <a className="text-gray-600 text-sm hover:text-purple-700">Company</a>
                            </li>
                            <li>
                                <a className="text-gray-600 text-sm hover:text-purple-700">News</a>
                            </li>
                            <li>
                                <a className="text-gray-600 text-sm hover:text-purple-700">Policies</a>
                            </li>
                            </nav>
                        </div>
                        <div className="lg:w-1/3 md:w-1/2 w-full px-4">
                            <h2 className="title-font font-medium text-purple-900 tracking-widest text-sm mb-3">FOLLOW US</h2>
                            <nav className="list-none mb-10">
                            <li>
                                <a className="text-gray-600 text-sm hover:text-purple-700">Facebook</a>
                            </li>
                            <li>
                                <a className="text-gray-600 text-sm hover:text-purple-700">Instagram</a>
                            </li>
                            <li>
                                <a className="text-gray-600 text-sm hover:text-purple-700">Twitter</a>
                            </li>
                            </nav>
                        </div>
                        <div className="lg:w-1/3 md:w-1/2 w-full pl-4">
                            <h2 className="title-font font-medium text-purple-900 tracking-widest text-sm mb-3">SUPPORT</h2>
                            <nav className="list-none mb-10">
                            <li>
                                <a className="text-gray-600 text-sm hover:text-purple-700">FAQ</a>
                            </li>
                            <li>
                                <a className="text-gray-600 text-sm hover:text-purple-700">Contacts</a>
                            </li>
                            </nav>
                        </div>
                    </div>
                </div>
                <div className=" container mx-auto flex-wrap justify-center w-screen py-2 px-5 flex max-w-full  bg-purple-100">
                        <p className="text-gray-600 text-xs ">© 2020 All Rights Reserved. EasySubmit</p>
                </div>
            </footer>
        </div>
    )
}

export default Footer;
