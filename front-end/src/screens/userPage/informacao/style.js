import { StyleSheet } from "react-native";

export default StyleSheet.create({
    busca:{
        alignSelf: "center",
        borderRadius:20,
        width:"75%"
    },
    buscaInput:{
        marginTop:80,
        backgroundColor:'gray',
        borderRadius:20,
        borderWidth:1,        
    },
    cardContent: {
        width: 70, margin: 10
    },
    cardContainer: {
        flexDirection: 'row',         
        alignItems: 'center',
        padding: 20,
        borderWidth: 2,
        width:"95%",
        alignSelf:"center",
        marginTop: 20
    },
    cardImage :{
        backgroundColor:"pink",
        width:100,
        height:100,
        resizeMode:"contain",
        borderRadius: 20
    },
    cardTitle:{
        marginLeft:10,
        fontWeight:"bold"
    },
    cardDescription:{
        flex:1,
        marginLeft:10
    }
})