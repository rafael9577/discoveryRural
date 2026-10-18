import { View, Text, TouchableOpacity, TextInput, Image, ActivityIndicator, Modal } from "react-native";
import { SafeAreaView, } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import Toast from "react-native-toast-message";
import { Feather } from "@expo/vector-icons"
import { useState } from 'react';
import styled from "./styled";

export default () => {

    return (
        <SafeAreaView style={styled.safeArea}>
            <View style= {styled.cardAnimal}>
                <Text>imagem</Text>
                <Text>Titulo</Text>
                <Text>Descrição</Text>
            </View>        
        </SafeAreaView>

)};