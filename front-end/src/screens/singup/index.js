import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, Modal } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import Toast from "react-native-toast-message";
import { Feather } from "@expo/vector-icons"
import { useState } from 'react';
import styled from "./styled";


export default () => {
    // blabla
    // navegação
    const navigation = useNavigation();

    // estados de memoria
    const [emailFild, setEmailFild] = useState("");
    const [passowrdFild, setPassowrdFild] = useState("");
    const [passowrdFild1, setPassowrdFild1] = useState("");
    const [seePass, setSeePass] = useState(true);
    const [seePass1, setSeePass1] = useState(true);

    // ações
    const handleSignUp = () => {
    // PRECISA TAMBEM CHAMAR A API PARA VALIDAR SE O USUARIO EXISTE?
    if (emailFild === '' || passowrdFild === '' || passowrdFild1 === '') { // Verifica se algum campo está vazio
        alert("Preencha todos os campos!");
        return; 
    }

    if (passowrdFild !== passowrdFild1) { // 2. Verifica se as senhas não conferem
        alert("As senhas não coincidem!");
        return; 
    }
        console.log("Botão 'Cadastrar' foi clicado. Cadastrado:", emailFild);};
    const toLogin = () => {
        console.log("Botão 'Já possui uma conta?' foi clicado");
        navigation.navigate("login");};


        

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
                        secureTextEntry={seePass} // para ocultar os caracteres
                        value={passowrdFild}   //
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
                        secureTextEntry={seePass1} // para ocultar os caracteres
                        value={passowrdFild1}   //
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
