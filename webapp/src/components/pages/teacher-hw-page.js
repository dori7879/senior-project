import Footer from '../footer';
import Header from '../header';
import React from 'react';
import axios from 'axios';

class TeacherHwPage extends React.Component{
    
    constructor(props) {
        super(props);
        this.onChangeGrade = this.onChangeGrade.bind(this);

        this.state = {
            grade: "",
            comment: "",
            title: "",
            description: "",
            files: [],
            closeDate: new Date()
        };          
    }
    onChangeGrade(e){
        this.setState({
            grade: e.target.value
        });
    }

    onChangeComment(e){
        this.setState({
            comment: e.target.value
        }); 
    }

    componentDidMount () {
        const { randomStr } = this.props.match.params
        axios.get(`/api/v1/homework-page/teacher/${randomStr}`)
          .then((response) => {
            this.setState(() => ({ 
                title: response.title,
                description: response.content,
                closeDate: response.closed_at
            }))
          })
      }
    render(){
        const submission = {
            student_fullName: "Dariya Shakenova",
            submission_time: Date(),
            content: "dhjfs,",
            attachments: [],
            comment:"ddcxv vkdcv ckvmv",
            grade: "A"
        }
        const isEmptyContent = submission.content.trim() === "";
        const isEmptyAttachment = submission.attachments.length === 0;
        const isGraded = !submission.grade.trim() === "";
        const isEmptyDesc = this.state.description.trim() === "";
        const isEmptyFile = this.state.files.length === 0; 
        return(
            <div>
                <Header />
                <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                    <div className="flex justify-center flex-col items-center pb-4">
                        <div className="text-purple-900 font-bold text-xl pt-4">Homework Submissions</div>
                    </div>
                        <form className="w-3/4 border border-purple-300 rounded bg-purple-300 p-4">
                            <div className=" flex flex-row"> 
                                <div className="flex flex-col">
                                    <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                                        <strong>Homework Title:</strong> <span className="text-purple-900">{this.state.title}</span></h2> 
                                    {
                                        isEmptyDesc ? null :
                                        <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                                            <strong>Description:</strong><br></br> <span className="text-purple-900">{this.state.description}</span></h2>
                                    }
                                    {
                                        isEmptyFile ? null :
                                        <h1> Attachments</h1> 
                                    } 
                                    <div className="border border-purple-700 rounded flex flex-col">
                                        <div className="flex flex-col items  pb-2">
                                            <p className="block tracking-wide text-gray-700 text-xs  px-4 pt-1">
                                                <strong>Student's Name:</strong> <span className="text-purple-900">{submission.student_fullName}</span></p>
                                            <p className="block tracking-wide text-gray-700 text-xs mb-2 px-4 pt-1">
                                                <strong>Submitted at:</strong> <span className="text-purple-900">{submission.submission_time}</span></p>  
                                            {
                                                isEmptyContent ? null :
                                                <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                    <strong>Content:</strong><br /> <span className="text-purple-900">{submission.content}</span></p>
                                            }
                                            {
                                                isEmptyAttachment ? null :
                                                <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                    <strong>Attachments:</strong><span className="text-purple-900">{submission.attachments}</span></p>
                                            }       
                                            {
                                                isGraded ?  <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                <strong>Grade:</strong><span className="text-purple-900">{submission.grade}</span>
                                                <strong>Comment:</strong><span className="text-purple-900">{submission.comment}</span>
                                                </p> :
                                                <div className="flex flex-col">
                                                    <div className="flex flex-row mb-2">
                                                        <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                            <strong>Grade: </strong></p>
                                                        <input  onChange={this.onChangeGrade} className=" w-10 text-gray-700 border border-purple-400 rounded text-xs py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="" />
                                                    </div>
                                                    <div className="flex flex-row">
                                                        <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                            <strong>Comment: </strong></p>
                                                        <textarea  onChange={this.onChangeComment} className=" text-gray-700 border border-purple-400 rounded text-xs py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="" />
                                                    </div>
                                                    <div>
                                                        <button type="submit" className="mb-2 ml-4 mt-2
                                                        relative flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                                            Grade Homework
                                                        </button>
                                                    </div>
                                                </div>
                                            }
                                        </div>
                                    </div>
                                </div>
                                <div className="flex flex-col ">
                                    <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                        <strong><span className="text-purple-800 font-bold">4</span> Homeworks Submitted</strong> </p>
                                    <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                        <strong></strong> </p>
                                    <button type="submit" className="mb-2 ml-4 mt-2
                                                relative flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                        Close Homework
                                    </button>
                                </div>
                            </div>
                        </form>
                </div>   
                <Footer />
            </div>
        )
    }
    }

export default TeacherHwPage;