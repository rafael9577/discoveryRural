import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, Modal } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import Toast from "react-native-toast-message";
import { Feather } from "@expo/vector-icons"
import { useState } from 'react';
import styled from "./styled";


export default () => {
    const navigation = useNavigation();


    const [emailFild, setEmailFild] = useState("");
    const [passowrdFild, setPassowrdFild] = useState("");
    const [passowrdFild1, setPassowrdFild1] = useState("");
    const [seePass, setSeePass] = useState(true);
    const [seePass1, setSeePass1] = useState(true);

    const handleSignUp = () => {
    // PRECISA TAMBEM CHAMAR A API PARA VALIDAR SE O USUARIO EXISTE?
    if (emailFild === '' || passowrdFild === '' || passowrdFild1 === '') {
        alert("Preencha todos os campos!");
        return; 
    }

    if (passowrdFild !== passowrdFild1) {
        alert("As senhas não coincidem!");
        return; 
    }
    }

    const toLogin = () => {
        navigation.navigate("login");
    };
    

        

    return (
        <SafeAreaView style={styled.safeArea}>
            <View style={styled.inputArea}>
                <View style={styled.singInPut}>
                    <TextInput
                        style={styled.textInput}
                        placeholder="insira seu e-mail"
                        placeholderTextColor="#666666"
                        value={emailFild}
                        onChangeText={t => setEmailFild(t)}
                    />
                </View>

                <View style={[styled.singInPut, { flexDirection: 'row', alignItems: 'center' }]}>
                    <TextInput
                        style={styled.textInput}
                        placeholder="insira a senha"
                        placeholderTextColor="#666666"
                        secureTextEntry={seePass} 
                        value={passowrdFild}   
                        onChangeText={t => setPassowrdFild(t)}
                    />
                    <TouchableOpacity onPress={() => setSeePass(!seePass)} style={{ marginRight: 12 }}>
                        <Feather 
                            name={seePass ? "eye-off" : "eye"} 
                            size={20} 
                            color="black" 
                        />
                    </TouchableOpacity>
                </View>

                <View style={[styled.singInPut, { flexDirection: 'row', alignItems: 'center' }]}>
                    <TextInput
                        style={styled.textInput}
                        placeholder="confirme a senha"
                        placeholderTextColor="#666666"
                        secureTextEntry={seePass1} 
                        value={passowrdFild1}   
                        onChangeText={t => setPassowrdFild1(t)}
                    />
                    <TouchableOpacity onPress={() => setSeePass1(!seePass1)} style={{ marginRight: 12 }}>
                        <Feather 
                            name={seePass1 ? "eye-off" : "eye"} 
                            size={20} 
                            color="black" 
                        />
                    </TouchableOpacity>
                </View>

            </View>
            <TouchableOpacity style={styled.custtomButton} onPress={handleSignUp}>
                <Text style={styled.customButtontext}>Cadastrar</Text>
            </TouchableOpacity>  

            <TouchableOpacity onPress={toLogin}>
                <Text style={styled.customButtontext2}> Já possui uma conta? </Text>
            </TouchableOpacity>
        </SafeAreaView>
    );
}
