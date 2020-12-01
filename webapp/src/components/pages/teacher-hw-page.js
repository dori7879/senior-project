import Footer from '../footer';
import Header from '../header';
import React from 'react';
import axios from 'axios';
import homework from '../../reducer/homework';

class TeacherHwPage extends React.Component{
    
    constructor(props) {
        super(props);
        this.onChangeGrade = this.onChangeGrade.bind(this);
        this.onChangeComments = this.onChangeComments.bind(this);
        this.handleGrade = this.handleGrade.bind(this);
        this.ckEditorRemoveTags = this.ckEditorRemoveTags.bind(this);

        this.state = {
            grade: "",
            comments: "",
            title: "",
            description: "",
            files: [],
            closeDate: new Date(),
            homeworks: []
        };          
    }
    onChangeGrade(e){
        this.setState({
            grade: e.target.value,
            comments: "",
            isGraded: false
        });
    }

    onChangeComments(e){
        this.setState({
            comments: e.target.value
        }); 
    }
    ckEditorRemoveTags (data) {     
        const editedData = data.replace('<p>', '').replace('</p>', '') 
        return editedData;
    }
    handleGrade(e){
        e.preventDefault();
        this.setState({
            isGraded: true
        });
        const { dispatch} = this.props;
       /* dispatch(gradeHomework( this.state.grade, this.state.comments))
            .then(() => {
                this.setState({
                    successful: true
                });
            })
            .catch(() => {
                this.setState({
                    successful: false
                });
            });*/
    }

    componentDidMount () {
        const { randomStr } = this.props.match.params;
        axios.get(`/api/v1/homework-page/teacher/${randomStr}`)
            .then((response) => {
                if (response.data) {
                    this.setState({
                        title: response.data.title,
                        description: response.data.content,
                        closeDate: response.data.closed_at,
                        homeworks: response.data.homeworks
                    })
                }
            })
      }
    render(){
      /*  const submission = {
            student_fullName: "Dariya Shakenova",
            submission_time: Date(),
            content: "dhjfs,",
            attachments: [],
            comments:"ddcxv vkdcv ckvmv",
            grade: "A"
        }*/
        console.log(this.state.homeworks);
        const data = this.ckEditorRemoveTags(this.state.description);
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
                                            <strong>Description:</strong><br></br> <span className="text-purple-900">{data}</span></h2>
                                    }
                                    {
                                        isEmptyFile ? null :
                                        <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                                            <strong>Attachments:</strong><br></br> <span className="text-purple-900">{this.state.files}</span></h2>
                                    } 
                                    {
                                        this.state.homeworks.map((homework, index) => (
                                            <div key={index} className="border border-purple-700 rounded flex flex-col">
                                                <div className="flex flex-col items  pb-2">
                                                    <p className="block tracking-wide text-gray-700 text-xs  px-4 pt-1">
                                                        <strong>Student's Name:</strong> <span className="text-purple-900">{homework.student_fullName}</span></p>
                                                    <p className="block tracking-wide text-gray-700 text-xs mb-2 px-4 pt-1">
                                                        <strong>Submitted at:</strong> <span className="text-purple-900">{homework.submitted_at}</span></p>  
                                                    {
                                                        homework.content.trim() === "" ? null :
                                                        <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                            <strong>Content:</strong><br /> <span className="text-purple-900">{homework.content}</span></p>
                                                    }
                                                    {
                                                        homework.attachments.length === 0 ? null :
                                                        <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                            <strong>Attachments:</strong><span className="text-purple-900">{homework.attachments}</span></p>
                                                    }       
                                                    {
                                                        homework.grade.trim()=== "" ?  <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                        <strong>Grade:</strong><span className="text-purple-900">{homework.grade}</span>
                                                        <strong>Comments:</strong><span className="text-purple-900">{homework.comments}</span>
                                                        </p> :
                                                        <div className="flex flex-col">
                                                            <div className="flex flex-row mb-2">
                                                                <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                                    <strong>Grade: </strong></p>
                                                                <input  onChange={this.onChangeGrade} className=" w-10 text-gray-700 border border-purple-400 rounded text-xs py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="" />
                                                            </div>
                                                            <div className="flex flex-row">
                                                                <p className="block tracking-wide text-gray-700 text-xs px-4 pt-1">
                                                                    <strong>Comments: </strong></p>
                                                                <textarea  onChange={this.onChangeComments} className=" text-gray-700 border border-purple-400 rounded text-xs py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="" />
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
                                        ))
                                    }
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